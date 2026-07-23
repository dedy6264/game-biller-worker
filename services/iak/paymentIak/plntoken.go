package paymentiak

import (
	"encoding/json"
	"fmt"
	"game-biller-worker/configs"
	"game-biller-worker/constans"
	"game-biller-worker/helpers"
	"game-biller-worker/models"
	"game-biller-worker/utils"
	"regexp"
	"strings"
	"time"
)

// PlnToken menangani payment PLN Token (prepaid)
// ProductReferenceID = 10
// additional = RefID
func PlnToken(request models.RequestPayment) (models.PaymentResult, error) {
	var (
		detailByte    []byte
		otherCustInfo models.PlnTokenBillDesc
	)

	// Ambil data customer dari BillDesc inquiry sebelumnya
	_ = json.Unmarshal([]byte(request.BillDesc), &otherCustInfo)

	baseURL := constans.IAK_DEV_BASE_URL
	if configs.APP_ENV != "DEV" {
		baseURL = constans.IAK_PROD_BASE_URL
	}
	apiURL := baseURL + constans.IAK_TOPUP_ENDPOINT

	// sign: md5(username + api_key + RefID)
	sign := helpers.SignIakEncrypt(request.RefID)

	username := constans.IAK_DEV_USERNAME
	if configs.APP_ENV != "DEV" {
		username = constans.IAK_PROD_USERNAME
	}

	// Assign data ke models.ReqPaymentPrepaidIAK
	iakRequest := models.ReqPaymentPrepaidIAK{
		CustomerId:  otherCustInfo.CustomerID,
		ProductCode: request.DataProduct.ProductCode,
		RefId:       request.RefID,
		Username:    username,
		Sign:        sign,
	}

	// Kirim request POST ke IAK
	respBytes, _, err := utils.WorkerRequestPOST(
		"json",
		apiURL,
		iakRequest,
		models.ReqHeader{},
		30*time.Second,
	)
	if err != nil {
		return models.PaymentResult{}, fmt.Errorf("failed to request IAK PLN token: %w", err)
	}

	// Handle response dengan fallback logic
	var iakResp models.RespPaymentPrepaidIAK
	if err := json.Unmarshal(respBytes, &iakResp); err != nil {
		return models.PaymentResult{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Fallback: jika Data.RefID == "" -> RespWorkerIakUndefined
	if iakResp.Data.RefID == "" {
		var undefined models.RespWorkerIakUndefined
		if err := json.Unmarshal(respBytes, &undefined); err == nil && undefined.ResponseCode != "" {
			iakResp.Data.Rc = undefined.ResponseCode
			iakResp.Data.Message = undefined.Message
		} else {
			var undefinedI models.RespWorkerIakUndefinedI
			if err := json.Unmarshal(respBytes, &undefinedI); err == nil && undefinedI.Data.ResponseCode != "" {
				iakResp.Data.Rc = undefinedI.Data.ResponseCode
				iakResp.Data.Message = undefinedI.Data.Message
			} else {
				var undefinedII models.RespWorkerIakUndefinedII
				if err := json.Unmarshal(respBytes, &undefinedII); err == nil {
					iakResp.Data.Rc = undefinedII.Data.Rc
					iakResp.Data.Message = undefinedII.Data.Message
				}
			}
		}
	}

	// Konversi response
	providerPayload := helpers.ProviderPayload{
		ResponseCode: iakResp.Data.Rc,
		Message:      iakResp.Data.Message,
		Case:         "PAY",
	}
	converted := helpers.IAKConverterResponse(providerPayload, 0)

	// Build BillDesc & parse SN jika sukses
	if converted.Code == helpers.CodeSuccess {
		iakResp.Data.Sn = strings.ReplaceAll(iakResp.Data.Sn, " ", "=")
		re := regexp.MustCompile(`[^a-zA-Z0-9-.=]+`)
		cleaned := re.ReplaceAllString(iakResp.Data.Sn, "/")
		dataCleaned := strings.Split(cleaned, "/")

		tarif, daya, kwh := "", "", ""
		if len(dataCleaned) > 2 {
			tarif = strings.ReplaceAll(dataCleaned[2], "=", "")
		}
		if len(dataCleaned) > 3 {
			daya = strings.ReplaceAll(dataCleaned[3], "=", "")
		}
		if len(dataCleaned) > 4 {
			kwh = strings.ReplaceAll(dataCleaned[4], "=", "")
		}
		if len(dataCleaned) > 0 {
			iakResp.Data.Sn = strings.ReplaceAll(dataCleaned[0], "=", "")
		}

		detail := models.PlnTokenBillDesc{
			CustomerID:   iakResp.Data.CustomerID,
			CustomerName: otherCustInfo.CustomerName,
			MeterNo:      otherCustInfo.MeterNo,
			Tarif:        tarif,
			Daya:         daya,
			Kwh:          kwh,
		}
		detailByte, _ = json.Marshal(detail)
	}

	result := models.PaymentResult{
		StatusCode:    converted.Code,
		RefID:         request.RefID,
		ProviderRefID: iakResp.Data.RefID,
		DataTransaction: models.DataTransaction{
			CustomerID:   iakResp.Data.CustomerID,
			SerialNumber: iakResp.Data.Sn,
			Price:        float64(iakResp.Data.Price),
			LastBalance:  float64(iakResp.Data.Balance),
		},
		ProviderDetail: models.ProviderFeedback{
			Code:    converted.CodeDetail,
			Message: converted.MessageDetail,
		},
		ProcessedAt: time.Now().Format(time.RFC3339),
		BillDesc:    string(detailByte),
	}
	return result, nil
}
