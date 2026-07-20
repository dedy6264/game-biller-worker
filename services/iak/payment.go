package iak

import (
	"encoding/json"
	"game-biller-worker/configs"
	"game-biller-worker/constans"
	"game-biller-worker/helpers"
	"game-biller-worker/models"
	"game-biller-worker/utils"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func Payment(c echo.Context) error {
	var (
		svc        = "IAK Payment"
		request    models.RequestPayment
		baseURL    = constans.IAK_DEV_BASE_URL
		detailByte []byte
	)

	// 1. Bind request body to models.RequestPayment
	if err := c.Bind(&request); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusBadRequest, models.PaymentResult{})
	}

	// 2. Determine base URL from environment

	if configs.APP_ENV != "DEV" {
		baseURL = constans.IAK_PROD_BASE_URL
	}
	apiURL := baseURL + constans.IAK_TOPUP_ENDPOINT

	// 3. Build sign: md5(username + api_key + additional)
	//    additional = reference number (RefID)
	sign := helpers.SignIakEncrypt(request.RefID)

	// 4. Assign data ke models.ReqPaymentPrepaidIak
	iakRequest := models.ReqPaymentPrepaidIAK{
		CustomerId:  request.ProviderRefID,
		ProductCode: "",
		RefId:       request.RefID,
		Username:    constans.IAK_DEV_USERNAME,
		Sign:        sign,
	}
	if configs.APP_ENV != "DEV" {
		iakRequest.Username = constans.IAK_PROD_USERNAME
	}

	// 5. Kirim request POST ke IAK
	respBytes, _, err := utils.WorkerRequestPOST(
		"json",
		apiURL,
		iakRequest,
		models.ReqHeader{},
		30*time.Second,
	)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to request")
		return c.JSON(http.StatusInternalServerError, models.PaymentResult{})
	}

	// 6. Handle response dengan fallback logic
	var iakResp models.RespPaymentPrepaidIAK

	// Try unmarshal default
	if err := json.Unmarshal(respBytes, &iakResp); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind response")
		return c.JSON(http.StatusInternalServerError, models.PaymentResult{})
	}

	// Fallback: jika Data.RefID == "" -> RespWorkerIakUndefined
	if iakResp.Data.RefID == "" {
		var undefined models.RespWorkerIakUndefined
		if err := json.Unmarshal(respBytes, &undefined); err == nil && undefined.ResponseCode != "" {
			iakResp.Data.Rc = undefined.ResponseCode
			iakResp.Data.Message = undefined.Message
		} else {
			// Fallback: jika ResponseCode == "" -> RespWorkerIakUndefinedI
			var undefinedI models.RespWorkerIakUndefinedI
			if err := json.Unmarshal(respBytes, &undefinedI); err == nil && undefinedI.Data.ResponseCode != "" {
				iakResp.Data.Rc = undefinedI.Data.ResponseCode
				iakResp.Data.Message = undefinedI.Data.Message
			} else {
				// Fallback: jika Data.ResponseCode == "" -> RespWorkerIakUndefinedII
				var undefinedII models.RespWorkerIakUndefinedII
				if err := json.Unmarshal(respBytes, &undefinedII); err == nil {
					iakResp.Data.Rc = undefinedII.Data.Rc
					iakResp.Data.Message = undefinedII.Data.Message
				}
			}
		}
	}
	// 7. Konversi response menggunakan helpers.IAKConverterResponse
	providerPayload := helpers.ProviderPayload{
		ResponseCode: iakResp.Data.Rc,
		Message:      iakResp.Data.Message,
		Case:         "PAY",
	}
	converted := helpers.IAKConverterResponse(providerPayload, 0)
	if converted.Code == "INQ-SYS-001" {
		switch request.DataProduct.ProductReferenceID {
		case 10: //plntoken
			var otherCustInfo models.PlnTokenBillDesc
			_ = json.Unmarshal([]byte(request.BillDesc), &otherCustInfo)
			iakResp.Data.Sn = strings.ReplaceAll(iakResp.Data.Sn, " ", "=")
			re := regexp.MustCompile(`[^a-zA-Z0-9-.=]+`)
			// Ganti semua karakter non-alphanumerik dengan string kosong
			cleaned := re.ReplaceAllString(iakResp.Data.Sn, "/")
			dataCleaned := strings.Split(cleaned, "/")
			tarif := strings.ReplaceAll(dataCleaned[2], "=", "")
			daya := strings.ReplaceAll(dataCleaned[3], "=", "")
			kwh := strings.ReplaceAll(dataCleaned[4], "=", "")
			iakResp.Data.Sn = strings.ReplaceAll(dataCleaned[0], "=", "")
			detail := models.PlnTokenBillDesc{
				CustomerID:   iakResp.Data.CustomerID,
				CustomerName: otherCustInfo.CustomerName,
				MeterNo:      otherCustInfo.MeterNo,
				Tarif:        tarif,
				Daya:         daya,
				Kwh:          kwh,
			}
			detailByte, _ = json.Marshal(detail)
		default:
			detail := models.BillDesc{
				CustomerID: iakResp.Data.CustomerID,
			}
			detailByte, _ = json.Marshal(detail)
		}
	}
	// 8. Return data baku dengan format models.PaymentResult
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
	return c.JSON(http.StatusOK, result)
}
