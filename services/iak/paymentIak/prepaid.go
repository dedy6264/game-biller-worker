package paymentiak

import (
	"encoding/json"
	"fmt"
	"game-biller-worker/configs"
	"game-biller-worker/constans"
	"game-biller-worker/helpers"
	"game-biller-worker/models"
	"game-biller-worker/utils"
	"time"
)

// Prepaid menangani payment produk prepaid umum (game topup, pulsa, dll)
// additional = RefID
func Prepaid(request models.RequestPayment) (models.PaymentResult, error) {
	var detailByte []byte

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
		CustomerId:  request.ProviderRefID,
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
		return models.PaymentResult{}, fmt.Errorf("failed to request IAK prepaid: %w", err)
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

	// Build BillDesc jika sukses
	if converted.Code == helpers.CodeSuccess {
		detail := models.BillDesc{
			CustomerID: iakResp.Data.CustomerID,
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
