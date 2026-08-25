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
			result.Code = CodePending
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
	"01":  {"INVOICE HAS BEEN PAID", "Payment request rejected. The referenced transaction is not in INQUIRY_SUCCESS status and cannot be processed.", CodeInvalidTransactionNoOrStatus},
	"02":  {"BILL UNPAID", "Payment request rejected. The referenced transaction is not in INQUIRY_SUCCESS status and cannot be processed.", CodeInvalidTransactionNoOrStatus},
	"04":  {"BILLING ID EXPIRED", "Payment request rejected. The referenced transaction is not in INQUIRY_SUCCESS status and cannot be processed.", CodeInvalidTransactionNoOrStatus},
	"06":  {"INQUIRY ID NOT FOUND", "External game validator service rejected the provided Player ID, Zone ID, or Server ID.", CodeInvalidCustID},
	"08":  {"BILLING ID BLOCKED", "External game validator service rejected the provided Player ID, Zone ID, or Server ID.", CodeInvalidCustID},
	"09":  {"INQUIRY FAILED", "The execution request was rejected by the distribution core gateway due to structural rule conflicts.", CodeServiceDisruption},
	"10":  {"BILL IS NOT AVAILABLE", "The requested target product code or SKU is missing or inactive in our catalog.", CodeInvalidCustID},
	"42":  {"PAYMENT REQUEST HAVEN'T BEEN RECEIVED", "Payment request rejected. The referenced transaction is not in INQUIRY_SUCCESS status and cannot be processed.", CodeInvalidTransactionNoOrStatus},
	"44":  {"EXCEEDING MAXIMAL DAILY INQUIRY ALLOWED", "The total value or volume of transactions has exceeded the agreed maximum daily limit.", CodeMaxTrx},
	"45":  {"TOO MANY INQUIRY REQUESTS", "The total value or volume of transactions has exceeded the agreed maximum daily limit.", CodeMaxTrx},
	"141": {"INVALID USER ID / ZONE ID / SERVER ID / ROLENAME", "External game validator service rejected the provided Player ID, Zone ID, or Server ID.", CodeInvalidCustID},
	"142": {"INVALID USER ID", "External game validator service rejected the provided Player ID, Zone ID, or Server ID.", CodeInvalidCustID},
	"143": {"INQUIRY NOT NEEDED", "Payment request rejected. The referenced transaction is not in INQUIRY_SUCCESS status and cannot be processed.", CodeInvalidTransactionNoOrStatus},

	// System / Merchant Error
	"07":  {"FAILED", "Your current transaction has failed. Please try again.", CodeServiceDisruption},
	"13":  {"CUSTOMER NUMBER BLOCKED", "", CodeInvalidCustID},
	"14":  {"INCORRECT DESTINATION NUMBER", "", CodeInvalidCustID},
	"16":  {"NUMBER NOT MATCH WITH OPERATOR", "", CodeInvalidCustID},
	"17":  {"INSUFFICIENT DEPOSIT", "", CodeServiceDisruption},
	"18":  {"NUMBER NOT AVAILABLE", "", CodeInvalidCustID},
	"19":  {"NUMBER IS ALREADY IN USE", "", CodeInvalidCustID},
	"20":  {"CODE NOT FOUND", "", CodeServiceDisruption},
	"21":  {"NUMBER EXPIRED", "", CodeInvalidCustID},
	"121": {"MONTHLY TOP UP LIMIT EXCEEDED", "", CodeMaxTrx},
	"131": {"TOP UP REGION BLOCKED FOR PLAYER", "", CodeInvalidCustID},
	"132": {"PRODUCT CODE NOT ELIGIBLE DUE TO SUBSCRIBER LOCATION", "", CodeInvalidCustID},
	"202": {"MAXIMUM 1 NUMBER 1 TIME IN 1 DAY", "", CodeMaxTrx},
	"203": {"NUMBER IS TOO LONG", "", CodeInvalidCustID},
	"206": {"THIS DESTINATION NUMBER HAS BEEN BLOCKED", "", CodeInvalidCustID},
	"207": {"MAXIMUM 1 NUMBER WITH ANY CODE 1 TIME IN 1 DAY", "", CodeMaxTrx},
	"301": {"EMAIL SEND LIMIT REACHED", "", CodeServiceDisruption},
	"91":  {"DATABASE CONNECTION ERROR", "Upstream provider responded with an unmapped error structure or critical raw payload exception.", CodeServiceDisruption},
	"92":  {"GENERAL ERROR", "Upstream provider responded with an unmapped error structure or critical raw payload exception.", CodeServiceDisruption},
	"93":  {"INVALID AMOUNT", "The destination identity string failed regex or basic character formatting check.", CodeServiceDisruption},
	"94":  {"SERVICE HAS EXPIRED", "The requested target product code or SKU is missing or inactive in our catalog.", CodeServiceDisruption},
	"100": {"INVALID SIGNATURE", "Invalid client key, invalid signature hash, or IP address not registered in the whitelist.", CodeServiceDisruption},
	"101": {"INVALID COMMAND", "The API endpoint requested does not support the HTTP verb or request method applied.", CodeServiceDisruption},
	"102": {"INVALID IP ADDRESS", "Invalid client key, invalid signature hash, or IP address not registered in the whitelist.", CodeServiceDisruption},
	"103": {"TIMEOUT", "Upstream server failed to respond within the designated execution window. Money safely refunded.", CodeServiceDisruption},
	"105": {"MISC ERROR / BILLER SYSTEM ERROR", "Upstream provider responded with an unmapped error structure or critical raw payload exception.", CodeServiceDisruption},
	"106": {"PRODUCT IS TEMPORARILY OUT OF SERVICE", "The upstream vendor gateway for this specific product is currently undergoing maintenance.", CodeServiceDisruption},
	"107": {"XML / JSON FORMAT ERROR", "The destination identity string failed regex or basic character formatting check.", CodeServiceDisruption},
	"110": {"SYSTEM UNDER MAINTENANCE", "The upstream vendor gateway for this specific product is currently undergoing maintenance.", CodeServiceDisruption},
	"117": {"PAGE NOT FOUND", "The API endpoint requested does not support the HTTP verb or request method applied.", CodeServiceDisruption},
	"204": {"WRONG AUTHENTICATION", "Invalid client key, invalid signature hash, or IP address not registered in the whitelist.", CodeServiceDisruption},
	"205": {"WRONG COMMAND", "The API endpoint requested does not support the HTTP verb or request method applied.", CodeServiceDisruption},
}

var respPay = map[string]mainResponse{
	// SUCCESS
	"00": {"SUCCESS / PAYMENT SUCCESS", "Payment confirmed and transaction successfully processed. Delivery is in progress.", CodeSuccess},

	// PENDING
	"05":  {"UNDEFINED ERROR", "Transaction accepted by core system and currently queued in the upstream gateway.", CodePending},
	"39":  {"PENDING / TRANSACTION IN PROCESS", "Transaction accepted by core system and currently queued in the upstream gateway.", CodePending},
	"91":  {"DATABASE CONNECTION ERROR", "Transaction accepted by core system and currently queued in the upstream gateway.", CodePending},
	"94":  {"SERVICE HAS EXPIRED", "Transaction accepted by core system and currently queued in the upstream gateway.", CodePending},
	"103": {"TIMEOUT", "Transaction accepted by core system and currently queued in the upstream gateway.", CodePending},
	"105": {"MISC ERROR / BILLER SYSTEM ERROR", "Transaction accepted by core system and currently queued in the upstream gateway.", CodePending},
	"110": {"SYSTEM UNDER MAINTENANCE", "Transaction accepted by core system and currently queued in the upstream gateway.", CodePending},
	"201": {"UNDEFINED RESPONSE CODE", "Transaction accepted by core system and currently queued in the upstream gateway.", CodePending},

	// FAILED
	"13":  {"CUSTOMER NUMBER BLOCKED", "", CodeInvalidCustID},
	"14":  {"INCORRECT DESTINATION NUMBER", "", CodeInvalidCustID},
	"16":  {"NUMBER NOT MATCH WITH OPERATOR", "", CodeInvalidCustID},
	"18":  {"NUMBER NOT AVAILABLE", "", CodeInvalidCustID},
	"19":  {"NUMBER IS ALREADY IN USE", "", CodeInvalidCustID},
	"20":  {"CODE NOT FOUND", "", CodeInvalidProduct},
	"21":  {"NUMBER EXPIRED", "", CodeInvalidCustID},
	"121": {"MONTHLY TOP UP LIMIT EXCEEDED", "", CodeMaxTrx},
	"131": {"TOP UP REGION BLOCKED FOR PLAYER", "Your current destination number top up request is blocked in that region. Please try again with a different destination number.", CodeInvalidCustID},
	"132": {"PRODUCT CODE NOT ELIGIBLE DUE TO SUBSCRIBER LOCATION", "", CodeInvalidProduct},
	"202": {"MAXIMUM 1 NUMBER 1 TIME IN 1 DAY", "", CodeMaxTrx},
	"203": {"NUMBER IS TOO LONG", "", CodeInvalidCustID},
	"206": {"THIS DESTINATION NUMBER HAS BEEN BLOCKED", "", CodeInvalidCustID},
	"207": {"MAXIMUM 1 NUMBER WITH ANY CODE 1 TIME IN 1 DAY", "", CodeMaxTrx},

	"17":  {"INSUFFICIENT DEPOSIT", "", CodeServiceDisruption},
	"301": {"EMAIL SEND LIMIT REACHED", "", CodeServiceDisruption},
	"07":  {"FAILED", "Your current transaction has failed. Please try again.", CodeServiceDisruption},
	"12":  {"BALANCE MAXIMUM LIMIT EXCEEDED", "", CodeServiceDisruption},
	"92":  {"GENERAL ERROR", "Upstream provider responded with an unmapped error structure or critical raw payload exception.", CodeServiceDisruption},
	"93":  {"INVALID AMOUNT", "The destination identity string failed regex or basic character formatting check.", CodeServiceDisruption},
	"100": {"INVALID SIGNATURE", "Invalid client key, invalid signature hash, or IP address not registered in the whitelist.", CodeServiceDisruption},
	"101": {"INVALID COMMAND", "The API endpoint requested does not support the HTTP verb or request method applied.", CodeServiceDisruption},
	"102": {"INVALID IP ADDRESS", "Invalid client key, invalid signature hash, or IP address not registered in the whitelist.", CodeServiceDisruption},
	"106": {"PRODUCT IS TEMPORARILY OUT OF SERVICE", "The upstream vendor gateway for this specific product is currently undergoing maintenance.", CodeServiceDisruption},
	"107": {"XML / JSON FORMAT ERROR", "The destination identity string failed regex or basic character formatting check.", CodeServiceDisruption},
	"117": {"PAGE NOT FOUND", "The API endpoint requested does not support the HTTP verb or request method applied.", CodeServiceDisruption},
	"204": {"WRONG AUTHENTICATION", "Invalid client key, invalid signature hash, or IP address not registered in the whitelist.", CodeServiceDisruption},
	"205": {"WRONG COMMAND", "The API endpoint requested does not support the HTTP verb or request method applied.", CodeServiceDisruption},
	"141": {"INVALID USER ID / ZONE ID / SERVER ID / ROLENAME", "External game validator service rejected the provided Player ID, Zone ID, or Server ID.", CodeInvalidCustID},
	"142": {"INVALID USER ID", "External game validator service rejected the provided Player ID, Zone ID, or Server ID.", CodeInvalidCustID},
	"143": {"INQUIRY NOT NEEDED", "Payment request rejected. The referenced transaction is not in INQUIRY_SUCCESS status and cannot be processed.", CodeServiceDisruption},
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
