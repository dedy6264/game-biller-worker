package inquiryiak

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

// PlnPrepaid menangani inquiry PLN Token (prepaid)
// CategoryID=4, TypeID=1
// additional = CustomerID
func PlnPrepaid(request models.RequestInquiry) (models.InquiryResult, error) {
	var (
		billdescByte []byte
	)

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

	respBytes, _, err := utils.WorkerRequestPOST(
		"json",
		apiURL,
		iakRequest,
		models.ReqHeader{},
		30*time.Second,
	)
	if err != nil {
		return models.InquiryResult{}, fmt.Errorf("failed to request IAK prepaid PLN: %w", err)
	}

	// Parse response ke RespInquiryPlnTokenIAK
	var iakResp models.RespInquiryPlnTokenIAK
	if err := json.Unmarshal(respBytes, &iakResp); err != nil {
		return models.InquiryResult{}, fmt.Errorf("failed to unmarshal response: %w", err)
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
	s, _ := json.Marshal(converted)
	fmt.Println("::", string(s))
	// Build BillDesc jika sukses
	if converted.Code == helpers.CodeInqSuccess {
		re := regexp.MustCompile(`[^a-zA-Z0-9-.]+`)
		cleaned := re.ReplaceAllString(iakResp.Data.SegmentPower, "/")
		segmentPowerArr := strings.Split(cleaned, "/")

		tarif, daya := "", ""
		if len(segmentPowerArr) > 0 {
			tarif = segmentPowerArr[0]
		}
		if len(segmentPowerArr) > 1 {
			daya = segmentPowerArr[1]
		}

		detail := models.PlnTokenBillDesc{
			CustomerID:   iakResp.Data.CustomerID,
			CustomerName: iakResp.Data.Name,
			MeterNo:      iakResp.Data.MeterNo,
			Tarif:        tarif,
			Daya:         daya,
		}
		billdescByte, _ = json.Marshal(detail)
	}

	result := models.InquiryResult{
		StatusCode: converted.Code,
		RefID:      request.RefID,
		DataTransaction: models.DataTransaction{
			CustomerID: iakResp.Data.CustomerID,
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
