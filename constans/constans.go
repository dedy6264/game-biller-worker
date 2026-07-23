package constans

const (
	// IAK Prepaid Dev Credentials
	IAK_DEV_USERNAME     = "082137789378"
	IAK_DEV_API_KEY      = "51562ac44252c544AhiD"
	IAK_DEV_BASE_URL     = "https://prepaid.iak.dev"
	IAK_DEV_POSTPAID_URL = "https://testpostpaid.mobilepulsa.net"

	// IAK Prepaid Prod Credentials
	IAK_PROD_USERNAME     = "082137789378"
	IAK_PROD_API_KEY      = "78362ac44e3786a5QWrQ"
	IAK_PROD_BASE_URL     = "https://prepaid.iak.id"
	IAK_PROD_POSTPAID_URL = "https://mobilepulsa.net"

	// IAK Endpoints
	IAK_TOPUP_ENDPOINT          = "/api/top-up"
	IAK_INQUIRY_PLN_ENDPOINT    = "/api/inquiry-pln"
	IAK_INQUIRY_POSTPAID_ENDPOINT = "/api/v1/bill/check"

	// Product Category ID
	PRODUCT_CATEGORY_PLN = int64(4)

	// Product Type ID
	PRODUCT_TYPE_PREPAID  = int64(1)
	PRODUCT_TYPE_POSTPAID = int64(2)

	// Product Reference ID
	PRODUCT_REFERENCE_PLN_TOKEN = int64(10)
)
