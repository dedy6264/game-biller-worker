package paymentiak

import (
	"encoding/json"
	"fmt"
	"game-biller-worker/configs"
	"game-biller-worker/constans"
	"game-biller-worker/helpers"
	"game-biller-worker/models"
	"game-biller-worker/utils"
	"strconv"
	"time"
)

// PlnToken menangani payment PLN Token (prepaid)
// ProductReferenceID = 10
// additional = RefID
func PlnPostpaid(request models.RequestPayment) (models.PaymentResult, error) {
	var (
		detailByte    []byte
		otherCustInfo models.PlnTokenBillDesc
		adminFee      float64
	)

	// Ambil data customer dari BillDesc inquiry sebelumnya
	_ = json.Unmarshal([]byte(request.BillDesc), &otherCustInfo)

	baseURL := constans.IAK_DEV_BASE_URL
	if configs.APP_ENV != "DEV" {
		baseURL = constans.IAK_PROD_BASE_URL
	}
	apiURL := baseURL + constans.IAK_INQUIRY_POSTPAID_ENDPOINT

	// sign: md5(username + api_key + RefID)
	sign := helpers.SignIakEncrypt(request.RefID)

	username := constans.IAK_DEV_USERNAME
	if configs.APP_ENV != "DEV" {
		username = constans.IAK_PROD_USERNAME
	}

	// Assign data ke models.ReqPaymentPrepaidIAK
	iakRequest := models.ReqPaymentPostpaidIAK{
		Commands: "pay-pasca",
		TrID:     request.RefID,
		Username: username,
		Sign:     sign,
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
		return models.PaymentResult{
			StatusCode: helpers.CodePending,
		}, fmt.Errorf("failed to request IAK PLN token: %w", err)
	}

	// Handle response dengan fallback logic
	var iakResp models.RespPaymentPlnIAK
	if err := json.Unmarshal(respBytes, &iakResp); err != nil {
		return models.PaymentResult{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Fallback: jika Data.RefID == "" -> RespWorkerIakUndefined
	if iakResp.Data.RefID == "" {
		var undefined models.RespWorkerIakUndefined
		if err := json.Unmarshal(respBytes, &undefined); err == nil && undefined.ResponseCode != "" {
			iakResp.Data.ResponseCode = undefined.ResponseCode
			iakResp.Data.Message = undefined.Message
		} else {
			var undefinedI models.RespWorkerIakUndefinedI
			if err := json.Unmarshal(respBytes, &undefinedI); err == nil && undefinedI.Data.ResponseCode != "" {
				iakResp.Data.ResponseCode = undefinedI.Data.ResponseCode
				iakResp.Data.Message = undefinedI.Data.Message
			} else {
				var undefinedII models.RespWorkerIakUndefinedII
				if err := json.Unmarshal(respBytes, &undefinedII); err == nil {
					iakResp.Data.ResponseCode = undefinedII.Data.Rc
					iakResp.Data.Message = undefinedII.Data.Message
				}
			}
		}
	}

	// Konversi response
	providerPayload := helpers.ProviderPayload{
		ResponseCode: iakResp.Data.ResponseCode,
		Message:      iakResp.Data.Message,
		Case:         "PAY",
	}
	converted := helpers.IAKConverterResponse(providerPayload, 0)

	// Build BillDesc & parse SN jika sukses
	if converted.Code == "SUC-INT-000" {
		adminFee = float64(iakResp.Data.Admin)
		var details []models.PlnDetail
		for _, data := range iakResp.Data.Desc.Tagihan.Detail {
			detail := models.PlnDetail{
				Periode:    data.Periode,
				Tagihan:    data.NilaiTagihan,
				Admin:      data.Admin,
				Denda:      data.Denda,
				MeterAwal:  data.MeterAwal,
				MeterAkhir: data.MeterAkhir,
			}
			details = append(details, detail)
		}
		// merchantFee = float64(iakResp.Data.Price - iakResp.Data.SellingPrice)
		billdesc := models.PlnTokenBillDesc{
			CustomerID:   iakResp.Data.Hp,
			CustomerName: iakResp.Data.TrName,
			Tarif:        iakResp.Data.Desc.Tarif,
			Daya:         strconv.Itoa(iakResp.Data.Desc.Daya),
			Details:      details,
			LembTag:      int64(len(iakResp.Data.Desc.Tagihan.Detail)),
		}
		detailByte, _ = json.Marshal(billdesc)
	}

	result := models.PaymentResult{
		StatusCode:    converted.Code,
		RefID:         request.RefID,
		ProviderRefID: iakResp.Data.RefID,
		DataTransaction: models.DataTransaction{
			CustomerID:   iakResp.Data.Hp,
			SerialNumber: iakResp.Data.Noref,
			Price:        float64(iakResp.Data.Price),
			LastBalance:  float64(iakResp.Data.Balance),
			AdminFee:     adminFee,
			GrandTotal:   float64(iakResp.Data.Price) + adminFee,
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
