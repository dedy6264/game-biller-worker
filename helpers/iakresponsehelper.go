package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"game-biller-worker/configs"
	"game-biller-worker/constans"
	"strings"
)

// ProviderPayload struct untuk menampung JSON response langsung dari provider (IAK)
type ProviderPayload struct {
	ResponseCode string `json:"rc"`
	Message      string `json:"message"`
	Description  string `json:"description"`
	Case         string `json:"case"`
}

// StandardizedResponse struct output hasil konversi akhir
type StandardizedResponse struct {
	CodeDetail        string `json:"code_detail"`        // RC Asli Provider
	MessageDetail     string `json:"message_detail"`     // Message Asli Provider
	DescriptionDetail string `json:"description_detail"` // Description Asli Provider
	Code              string `json:"code"`               // SUCCESS | PENDING | FAILED
}

func IAKConverterResponse(payload ProviderPayload, status int) (result StandardizedResponse) {
	result = StandardizedResponse{
		CodeDetail:        payload.ResponseCode,
		MessageDetail:     payload.Message,
		DescriptionDetail: payload.Description,
	}
	switch strings.ToUpper(payload.Case) {
	case "INQ":
		resp, ok := respInq[payload.ResponseCode]
		if !ok {
			result.Code = CodeErrPvd4305
			return
		}
		result.Code = resp.Maincode
		if result.MessageDetail == "" {
			result.MessageDetail = resp.ProviderMsg
		}
		if result.DescriptionDetail == "" {
			result.DescriptionDetail = resp.MainMsg
		}
	default:
		resp, ok := respPay[payload.ResponseCode]
		if !ok {
			result.Code = CodeIntrPending
			return
		}
		result.Code = resp.Maincode
		if result.MessageDetail == "" {
			result.MessageDetail = resp.ProviderMsg
		}
		if result.DescriptionDetail == "" {
			result.DescriptionDetail = resp.MainMsg
		}
	}
	return
}

type mainResponse struct {
	ProviderMsg string
	MainMsg     string
	Maincode    string
}

var respInq = map[string]mainResponse{
	// SUCCESS
	"00": {"INQUIRY SUCCESS", "Inquiry accepted. Product and pricing verified. Awaiting payment confirmation from merchant.", CodeInqSuccess},

	// ERR User Input / Data
	"01":  {"INVOICE HAS BEEN PAID", "Payment request rejected. The referenced transaction is not in INQUIRY_SUCCESS status and cannot be processed.", CodeInvalidTransaction},
	"02":  {"BILL UNPAID", "Payment request rejected. The referenced transaction is not in INQUIRY_SUCCESS status and cannot be processed.", CodeInvalidTransaction},
	"04":  {"BILLING ID EXPIRED", "Payment request rejected. The referenced transaction is not in INQUIRY_SUCCESS status and cannot be processed.", CodeInvalidTransaction},
	"06":  {"INQUIRY ID NOT FOUND", "External game validator service rejected the provided Player ID, Zone ID, or Server ID.", CodeInvalidIdGame},
	"08":  {"BILLING ID BLOCKED", "External game validator service rejected the provided Player ID, Zone ID, or Server ID.", CodeInvalidIdGame},
	"09":  {"INQUIRY FAILED", "The execution request was rejected by the distribution core gateway due to structural rule conflicts.", CodeErrPvd3303},
	"10":  {"BILL IS NOT AVAILABLE", "The requested target product code or SKU is missing or inactive in our catalog.", CodeInvalidProductNotFound},
	"42":  {"PAYMENT REQUEST HAVEN'T BEEN RECEIVED", "Payment request rejected. The referenced transaction is not in INQUIRY_SUCCESS status and cannot be processed.", CodeInvalidTransaction},
	"44":  {"EXCEEDING MAXIMAL DAILY INQUIRY ALLOWED", "The total value or volume of transactions has exceeded the agreed maximum daily limit.", CodeErrInt202},
	"45":  {"TOO MANY INQUIRY REQUESTS", "The total value or volume of transactions has exceeded the agreed maximum daily limit.", CodeErrInt202},
	"141": {"INVALID USER ID / ZONE ID / SERVER ID / ROLENAME", "External game validator service rejected the provided Player ID, Zone ID, or Server ID.", CodeInvalidIdGame},
	"142": {"INVALID USER ID", "External game validator service rejected the provided Player ID, Zone ID, or Server ID.", CodeInvalidIdGame},
	"143": {"INQUIRY NOT NEEDED", "Payment request rejected. The referenced transaction is not in INQUIRY_SUCCESS status and cannot be processed.", CodeInvalidTransaction},

	// System / Merchant Error
	"91":  {"DATABASE CONNECTION ERROR", "Upstream provider responded with an unmapped error structure or critical raw payload exception.", CodeErrPvd4305},
	"92":  {"GENERAL ERROR", "Upstream provider responded with an unmapped error structure or critical raw payload exception.", CodeErrPvd4305},
	"93":  {"INVALID AMOUNT", "The destination identity string failed regex or basic character formatting check.", CodeInvalidCustId},
	"94":  {"SERVICE HAS EXPIRED", "The requested target product code or SKU is missing or inactive in our catalog.", CodeInvalidProductNotFound},
	"100": {"INVALID SIGNATURE", "Invalid client key, invalid signature hash, or IP address not registered in the whitelist.", CodeErrInt200},
	"101": {"INVALID COMMAND", "The API endpoint requested does not support the HTTP verb or request method applied.", CodeErrInt203},
	"102": {"INVALID IP ADDRESS", "Invalid client key, invalid signature hash, or IP address not registered in the whitelist.", CodeErrInt200},
	"103": {"TIMEOUT", "Upstream server failed to respond within the designated execution window. Money safely refunded.", CodeErrPvd1302},
	"105": {"MISC ERROR / BILLER SYSTEM ERROR", "Upstream provider responded with an unmapped error structure or critical raw payload exception.", CodeErrPvd4305},
	"106": {"PRODUCT IS TEMPORARILY OUT OF SERVICE", "The upstream vendor gateway for this specific product is currently undergoing maintenance.", CodeErrPvd1300},
	"107": {"XML / JSON FORMAT ERROR", "The destination identity string failed regex or basic character formatting check.", CodeInvalidCustId},
	"110": {"SYSTEM UNDER MAINTENANCE", "The upstream vendor gateway for this specific product is currently undergoing maintenance.", CodeErrPvd1300},
	"117": {"PAGE NOT FOUND", "The API endpoint requested does not support the HTTP verb or request method applied.", CodeErrInt203},
	"204": {"WRONG AUTHENTICATION", "Invalid client key, invalid signature hash, or IP address not registered in the whitelist.", CodeErrInt200},
	"205": {"WRONG COMMAND", "The API endpoint requested does not support the HTTP verb or request method applied.", CodeErrInt203},
}

var respPay = map[string]mainResponse{
	// SUCCESS
	"00": {"SUCCESS / PAYMENT SUCCESS", "Payment confirmed and transaction successfully processed. Delivery is in progress.", CodeSuccess},

	// PENDING
	"05":  {"UNDEFINED ERROR", "Transaction accepted by core system and currently queued in the upstream gateway.", CodeIntrPending},
	"39":  {"PENDING / TRANSACTION IN PROCESS", "Transaction accepted by core system and currently queued in the upstream gateway.", CodeIntrPending},
	"91":  {"DATABASE CONNECTION ERROR", "Transaction accepted by core system and currently queued in the upstream gateway.", CodeIntrPending},
	"94":  {"SERVICE HAS EXPIRED", "Transaction accepted by core system and currently queued in the upstream gateway.", CodeIntrPending},
	"103": {"TIMEOUT", "Transaction accepted by core system and currently queued in the upstream gateway.", CodeIntrPending},
	"105": {"MISC ERROR / BILLER SYSTEM ERROR", "Transaction accepted by core system and currently queued in the upstream gateway.", CodeIntrPending},
	"110": {"SYSTEM UNDER MAINTENANCE", "Transaction accepted by core system and currently queued in the upstream gateway.", CodeIntrPending},
	"201": {"UNDEFINED RESPONSE CODE", "Transaction accepted by core system and currently queued in the upstream gateway.", CodeIntrPending},

	// FAILED
	"92":  {"GENERAL ERROR", "Upstream provider responded with an unmapped error structure or critical raw payload exception.", CodeErrPvd4305},
	"93":  {"INVALID AMOUNT", "The destination identity string failed regex or basic character formatting check.", CodeInvalidCustId},
	"100": {"INVALID SIGNATURE", "Invalid client key, invalid signature hash, or IP address not registered in the whitelist.", CodeErrInt200},
	"101": {"INVALID COMMAND", "The API endpoint requested does not support the HTTP verb or request method applied.", CodeErrInt203},
	"102": {"INVALID IP ADDRESS", "Invalid client key, invalid signature hash, or IP address not registered in the whitelist.", CodeErrInt200},
	"106": {"PRODUCT IS TEMPORARILY OUT OF SERVICE", "The upstream vendor gateway for this specific product is currently undergoing maintenance.", CodeErrPvd1300},
	"107": {"XML / JSON FORMAT ERROR", "The destination identity string failed regex or basic character formatting check.", CodeInvalidCustId},
	"117": {"PAGE NOT FOUND", "The API endpoint requested does not support the HTTP verb or request method applied.", CodeErrInt203},
	"204": {"WRONG AUTHENTICATION", "Invalid client key, invalid signature hash, or IP address not registered in the whitelist.", CodeErrInt200},
	"205": {"WRONG COMMAND", "The API endpoint requested does not support the HTTP verb or request method applied.", CodeErrInt203},
}

func SignIakEncrypt(additional string) (sign string) {
	var (
		username, apikey string
	)
	if configs.APP_ENV == "DEV" {
		apikey = constans.IAK_DEV_API_KEY
		username = constans.IAK_DEV_USERNAME
	} else {
		apikey = constans.IAK_PROD_API_KEY
		username = constans.IAK_PROD_USERNAME
	}
	// sign: md5({username}+{api_key}+{additional})
	key := username + apikey + additional
	sign = createHash(key)
	return
}
func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
