package inquiryiak

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

// PlnPostpaid menangani inquiry PLN Pascabayar (postpaid)
// CategoryID=4, TypeID=2
// additional = RefID
func PlnPostpaid(request models.RequestInquiry) (models.InquiryResult, error) {
	var (
		billdescByte []byte
	)

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

	respBytes, _, err := utils.WorkerRequestPOST(
		"json",
		apiURL,
		iakRequest,
		models.ReqHeader{},
		30*time.Second,
	)
	if err != nil {
		return models.InquiryResult{}, fmt.Errorf("failed to request IAK postpaid PLN: %w", err)
	}

	// Parse response ke RespInquiryPlnIAK
	var iakResp models.RespInquiryPlnIAK
	if err := json.Unmarshal(respBytes, &iakResp); err != nil {
		return models.InquiryResult{}, fmt.Errorf("failed to unmarshal response: %w", err)
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

	// Build BillDesc jika sukses
	if converted.Code == helpers.CodeInqSuccess {
		var details []models.PlnDetail
		for _, d := range iakResp.Data.Desc.Tagihan.Detail {
			details = append(details, models.PlnDetail{
				Periode: d.Periode,
				Tagihan: d.NilaiTagihan,
				Admin:   d.Admin,
				Denda:   d.Denda,
			})
		}
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
	return result, nil
}
