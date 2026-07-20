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
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func Inquiry(c echo.Context) error {
	var (
		svc          = "IAK Inquiry"
		request      models.RequestInquiry
		billdescByte []byte
		adminFee     float64
	)

	// 1. Bind request body ke models.RequestInquiry
	if err := c.Bind(&request); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusBadRequest, models.InquiryResult{})
	}

	isPLN := request.DataProduct.ProductCategoryID == constans.PRODUCT_CATEGORY_PLN
	isPrepaid := request.DataProduct.ProductTypeID == constans.PRODUCT_TYPE_PREPAID
	isPostpaid := request.DataProduct.ProductTypeID == constans.PRODUCT_TYPE_POSTPAID

	var (
		respBytes []byte
		apiErr    error
	)

	// ─────────────────────────────────────────────────────────────
	// CASE A: PLN Prepaid (CategoryID=4, TypeID=1)
	// endpoint: /api/inquiry-pln | additional = CustomerID
	// ─────────────────────────────────────────────────────────────
	if isPLN && isPrepaid {
		baseURL := constans.IAK_DEV_BASE_URL
		if configs.APP_ENV != "DEV" {
			baseURL = constans.IAK_PROD_BASE_URL
		}
		apiURL := baseURL + constans.IAK_INQUIRY_PLN_ENDPOINT

		// sign: md5(username + api_key + CustomerID)
		sign := helpers.SignIakEncrypt(request.CustomerID)

		username := constans.IAK_DEV_USERNAME
		if configs.APP_ENV != "DEV" {
			username = constans.IAK_PROD_USERNAME
		}

		// Assign ke models.ReqInquiryPlnTokenIAK
		iakRequest := models.ReqInquiryPlnTokenIAK{
			Username:   username,
			CustomerID: request.CustomerID,
			Sign:       sign,
		}

		respBytes, _, apiErr = utils.WorkerRequestPOST(
			"json",
			apiURL,
			iakRequest,
			models.ReqHeader{},
			30*time.Second,
		)
		if apiErr != nil {
			helpers.ProcessLogger(c, svc, apiErr.Error(), "Failed to request IAK prepaid PLN")
			return c.JSON(http.StatusInternalServerError, models.InquiryResult{})
		}

		// Parse response ke RespPaymentPrepaidIAK
		var iakResp models.RespInquiryPlnTokenIAK
		if err := json.Unmarshal(respBytes, &iakResp); err != nil {
			helpers.ProcessLogger(c, svc, err.Error(), "Failed to unmarshal response")
			return c.JSON(http.StatusInternalServerError, models.InquiryResult{})
		}

		// Fallback: jika Data.SubscriberID == "" -> RespWorkerIakUndefined
		if iakResp.Data.SubscriberID == "" {
			var undefined models.RespWorkerIakUndefined
			if err := json.Unmarshal(respBytes, &undefined); err == nil && undefined.ResponseCode != "" {
				iakResp.Data.Rc = undefined.ResponseCode
				iakResp.Data.Message = undefined.Message
			} else {
				// Fallback: jika ResponseCode == "" -> RespWorkerIakUndefinedI
				var undefinedI models.RespWorkerIakUndefinedI
				if err := json.Unmarshal(respBytes, &undefinedI); err == nil {
					iakResp.Data.Rc = undefinedI.Data.ResponseCode
					iakResp.Data.Message = undefinedI.Data.Message
				}
			}
		}

		// Konversi response
		providerPayload := helpers.ProviderPayload{
			ResponseCode: iakResp.Data.Rc,
			Message:      iakResp.Data.Message,
			Case:         "INQ",
		}
		converted := helpers.IAKConverterResponse(providerPayload, 0)
		if converted.Code == "INQ-SYS-001" {
			re := regexp.MustCompile(`[^a-zA-Z0-9-.]+`)
			// Ganti semua karakter non-alphanumerik dengan string kosong
			cleaned := re.ReplaceAllString(iakResp.Data.SegmentPower, "/")
			segmentPowerArr := strings.Split(cleaned, "/")
			detail := models.PlnTokenBillDesc{
				CustomerID:   iakResp.Data.CustomerID,
				CustomerName: iakResp.Data.Name,
				MeterNo:      iakResp.Data.MeterNo,
				Tarif:        segmentPowerArr[0],
				Daya:         segmentPowerArr[1],
			}
			billdescByte, _ = json.Marshal(detail)
		}
		result := models.InquiryResult{
			StatusCode: converted.Code,
			RefID:      request.RefID,
			// ProviderRefID: iakResp.Data.RefID,
			DataTransaction: models.DataTransaction{
				AdminFee:   adminFee,
				CustomerID: iakResp.Data.CustomerID,
			},
			ProviderDetail: models.ProviderFeedback{
				Code:    converted.CodeDetail,
				Message: converted.MessageDetail,
			},
			ProcessedAt: time.Now().Format(time.RFC3339),
			BillDesc:    string(billdescByte),
		}

		return c.JSON(http.StatusOK, result)
	}

	// ─────────────────────────────────────────────────────────────
	// CASE B: PLN Postpaid (CategoryID=4, TypeID=2)
	// endpoint: /api/v1/bill/check | additional = RefID
	// ─────────────────────────────────────────────────────────────
	if isPLN && isPostpaid {
		baseURL := constans.IAK_DEV_POSTPAID_URL
		if configs.APP_ENV != "DEV" {
			baseURL = constans.IAK_PROD_POSTPAID_URL
		}
		apiURL := baseURL + constans.IAK_INQUIRY_POSTPAID_ENDPOINT

		// sign: md5(username + api_key + RefID)
		sign := helpers.SignIakEncrypt(request.RefID)

		username := constans.IAK_DEV_USERNAME
		if configs.APP_ENV != "DEV" {
			username = constans.IAK_PROD_USERNAME
		}

		// Assign ke models.ReqInquiryPostpaidIAK
		iakRequest := models.ReqInquiryPostpaidIAK{
			Commands: "inq-pasca",
			Hp:       request.CustomerID,
			Code:     request.DataProduct.ProductCode,
			RefId:    request.RefID,
			Username: username,
			Sign:     sign,
		}

		respBytes, _, apiErr = utils.WorkerRequestPOST(
			"json",
			apiURL,
			iakRequest,
			models.ReqHeader{},
			30*time.Second,
		)
		if apiErr != nil {
			helpers.ProcessLogger(c, svc, apiErr.Error(), "Failed to request IAK postpaid PLN")
			return c.JSON(http.StatusInternalServerError, models.InquiryResult{})
		}

		// Parse response ke RespPaymentPostpaidIAK
		var iakResp models.RespInquiryPlnIAK
		if err := json.Unmarshal(respBytes, &iakResp); err != nil {
			helpers.ProcessLogger(c, svc, err.Error(), "Failed to unmarshal response")
			return c.JSON(http.StatusInternalServerError, models.InquiryResult{})
		}

		// Fallback: jika Data.RefID == "" -> RespWorkerIakUndefined
		if iakResp.Data.RefID == "" {
			var undefined models.RespWorkerIakUndefined
			if err := json.Unmarshal(respBytes, &undefined); err == nil && undefined.ResponseCode != "" {
				iakResp.Data.Code = undefined.ResponseCode
				iakResp.Data.Message = undefined.Message
			} else {
				// Fallback: jika ResponseCode == "" -> RespWorkerIakUndefinedI
				var undefinedI models.RespWorkerIakUndefinedI
				if err := json.Unmarshal(respBytes, &undefinedI); err == nil && undefinedI.Data.ResponseCode != "" {
					iakResp.Data.Code = undefinedI.Data.ResponseCode
					iakResp.Data.Message = undefinedI.Data.Message
				} else {
					// Fallback: jika Data.ResponseCode == "" -> RespWorkerIakUndefinedII
					var undefinedII models.RespWorkerIakUndefinedII
					if err := json.Unmarshal(respBytes, &undefinedII); err == nil {
						iakResp.Data.Code = undefinedII.Data.Rc
						iakResp.Data.Message = undefinedII.Data.Message
					}
				}
			}
		}

		// Konversi response
		providerPayload := helpers.ProviderPayload{
			ResponseCode: iakResp.Data.Code,
			Message:      iakResp.Data.Message,
			Case:         "INQ",
		}
		converted := helpers.IAKConverterResponse(providerPayload, 0)
		if converted.Code == "INQ-SYS-001" {
			adminFee = float64(iakResp.Data.Admin)
			var details []models.PlnDetail
			for _, data := range iakResp.Data.Desc.Tagihan.Detail {
				detail := models.PlnDetail{
					Periode: data.Periode,
					Tagihan: data.NilaiTagihan,
					Admin:   data.Admin,
					Denda:   data.Denda,
				}
				details = append(details, detail)
			}
			// merchantFee = float64(iakResp.Data.Price - iakResp.Data.SellingPrice)
			billdesc := models.PlnTokenBillDesc{
				CustomerID:   iakResp.Data.Hp,
				CustomerName: iakResp.Data.TrName,
				Tarif:        iakResp.Data.Desc.Tarif,
				Daya:         strconv.Itoa(iakResp.Data.Desc.Daya),
				LembTag:      int64(len(iakResp.Data.Desc.Tagihan.Detail)),
				Details:      details,
			}
			billdescByte, _ = json.Marshal(billdesc)
		}
		result := models.InquiryResult{
			StatusCode:    converted.Code,
			RefID:         request.RefID,
			ProviderRefID: iakResp.Data.RefID,
			DataTransaction: models.DataTransaction{
				CustomerID: iakResp.Data.Hp,
				Price:      float64(iakResp.Data.Price),
				AdminFee:   float64(iakResp.Data.Admin),
				GrandTotal: float64(iakResp.Data.Price) + float64(iakResp.Data.Admin),
			},
			ProviderDetail: models.ProviderFeedback{
				Code:    converted.CodeDetail,
				Message: converted.MessageDetail,
			},
			ProcessedAt: time.Now().Format(time.RFC3339),
			BillDesc:    string(billdescByte),
		}

		return c.JSON(http.StatusOK, result)
	}

	// Kategori/type tidak dikenali
	return c.JSON(http.StatusBadRequest, models.InquiryResult{
		StatusCode: "ERR-VAL-100",
	})
}
