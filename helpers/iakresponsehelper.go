package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"game-biller-worker/configs"
	"strings"
)

// Status Kategori Internal Sistem Anda

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
		result.Code = respInq[payload.ResponseCode].Maincode
		if result.Code == "" {
			result.Code = "ERR-PVD4-306"
			return
		}
	default:
		result.Code = respPay[payload.ResponseCode].Maincode
		if result.Code == "" {
			result.Code = "PEN-INT-002"
			return
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
	"00": {"INQUIRY SUCCESS", "Your inquiry is successfully processed.", "INQ-SYS-001"},

	// ERR User Input / Data
	"01":  {"INVOICE HAS BEEN PAID", "The invoice with your inputted data has already been paid.", "ERR-INT-102"},
	"02":  {"BILL UNPAID", "Your bill is unpaid, only reaching inquiry status.", "ERR-INT-102"},
	"04":  {"BILLING ID EXPIRED", "Your reference ID (ref_id) is expired. Please inquiry and payment on the same day.", "ERR-INT-102"},
	"06":  {"INQUIRY ID NOT FOUND", "The inquiry ID (tr_id) that you've inputted is not found.", "ERR-VAL-100"},
	"08":  {"BILLING ID BLOCKED", "The customer ID for your inputted product code is blocked by IAK.", "ERR-VAL-100"},
	"09":  {"INQUIRY FAILED", "Your inquiry process failed. Please try to do the inquiry again.", "ERR-PVD3-303"},
	"10":  {"BILL IS NOT AVAILABLE", "The bill isn't available yet. Please try again when the bill is already available.", "ERR-INT-101"},
	"42":  {"PAYMENT REQUEST HAVEN'T BEEN RECEIVED", "Your current transaction is still in the inquiry process.", "ERR-INT-102"},
	"44":  {"EXCEEDING MAXIMAL DAILY INQUIRY ALLOWED", "Your transaction is exceeding your today's inquiry limit (100 times).", "ERR-INT-202"},
	"45":  {"TOO MANY INQUIRY REQUESTS", "Your current request is failed due too much inquiry request.", "ERR-INT-202"},
	"141": {"INVALID USER ID / ZONE ID / SERVER ID / ROLENAME", "Your inputted user ID / Zone ID / Server ID / Role name isn't valid.", "ERR-VAL-100"},
	"142": {"INVALID USER ID", "Your current destination number (user id) top up request is invalid.", "ERR-VAL-100"},
	"143": {"INQUIRY NOT NEEDED", "The inputted operator is a voucher type therefore it doesn't need player id validation.", "ERR-INT-102"},

	// System / Merchant Error
	"91":  {"DATABASE CONNECTION ERROR", "Error on the database connection.", "ERR-PVD4-305"},
	"92":  {"GENERAL ERROR", "The received response code is undefined yet.", "ERR-PVD4-305"},
	"93":  {"INVALID AMOUNT", "The amount inputted isn't valid.", "ERR-VAL-104"},
	"94":  {"SERVICE HAS EXPIRED", "Service has expired.", "ERR-INT-101"},
	"100": {"INVALID SIGNATURE", "Sign field doesn't contain the right key.", "ERR-INT-200"},
	"101": {"INVALID COMMAND", "Command field is not a valid command.", "ERR-INT-203"},
	"102": {"INVALID IP ADDRESS", "IP address isn't allowed to make a transaction.", "ERR-INT-200"},
	"103": {"TIMEOUT", "Request exceeds the timeout limit.", "ERR-PVD1-302"},
	"105": {"MISC ERROR / BILLER SYSTEM ERROR", "Error from the supplier/biller.", "ERR-PVD4-305"},
	"106": {"PRODUCT IS TEMPORARILY OUT OF SERVICE", "Product is in non-active status.", "ERR-PVD1-300"},
	"107": {"XML / JSON FORMAT ERROR", "Request body format isn't correct.", "ERR-VAL-104"},
	"110": {"SYSTEM UNDER MAINTENANCE", "System is currently under maintenance.", "ERR-PVD1-300"},
	"117": {"PAGE NOT FOUND", "API URL target is not found.", "ERR-INT-203"},
	"204": {"WRONG AUTHENTICATION", "Sign field authentication error.", "ERR-INT-200"},
	"205": {"WRONG COMMAND", "Command field value error.", "ERR-INT-203"},
}

var respPay = map[string]mainResponse{
	// SUCCESS
	"00": {"SUCCESS / PAYMENT SUCCESS", "Your payment / top up is successfully processed.", "SUC-INT-000"},

	// PENDING
	"05":  {"UNDEFINED ERROR", "Your transaction pending because of an undefined error.", "PEN-SYS-001"},
	"39":  {"PENDING / TRANSACTION IN PROCESS", "Your current transaction is being processed, please wait until fully processed.", "PEN-SYS-001"},
	"91":  {"DATABASE CONNECTION ERROR", "Error on the database connection.", "PEN-SYS-003"},
	"94":  {"SERVICE HAS EXPIRED", "Service has expired.", "PEN-INT-002"},
	"103": {"TIMEOUT", "Request exceeds the timeout limit.", "PEN-SYS-003"},
	"105": {"MISC ERROR / BILLER SYSTEM ERROR", "Error from the supplier/biller.", "PEN-SYS-003"},
	"110": {"SYSTEM UNDER MAINTENANCE", "System is currently under maintenance.", "PEN-SYS-003"},
	"201": {"UNDEFINED RESPONSE CODE", "The received response code is undefined yet.", "PEN-SYS-001"},

	// FAILED
	"92":  {"GENERAL ERROR", "The received response code is undefined yet.", "ERR-PVD4-305"},
	"93":  {"INVALID AMOUNT", "The amount inputted isn't valid.", "ERR-VAL-104"},
	"100": {"INVALID SIGNATURE", "Sign field doesn't contain the right key.", "ERR-INT-200"},
	"101": {"INVALID COMMAND", "Command field is not a valid command.", "ERR-INT-203"},
	"102": {"INVALID IP ADDRESS", "IP address isn't allowed to make a transaction.", "ERR-INT-200"},
	"106": {"PRODUCT IS TEMPORARILY OUT OF SERVICE", "Product is in non-active status.", "ERR-PVD1-300"},
	"107": {"XML / JSON FORMAT ERROR", "Request body format isn't correct.", "ERR-VAL-104"},
	"117": {"PAGE NOT FOUND", "API URL target is not found.", "ERR-INT-203"},
	"204": {"WRONG AUTHENTICATION", "Sign field authentication error.", "ERR-INT-200"},
	"205": {"WRONG COMMAND", "Command field value error.", "ERR-INT-203"},
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
