package models

type TransactionStatus string

const (
	StatusSuccess TransactionStatus = "SUCCESS"
	StatusPending TransactionStatus = "PENDING"
	StatusFailed  TransactionStatus = "FAILED"
)

type (
	PaymentResult struct {
		CustomerID     string            `json:"customer_id"`
		RefID          string            `json:"ref_id"`          // Reference ID dari Main Service
		ProviderRefID  string            `json:"provider_ref_id"` // Transaction ID / Invoice ID dari Provider
		Status         TransactionStatus `json:"status"`          // SUCCESS | PENDING | FAILED- blm tau mau dipake atau ga
		SerialNumber   string            `json:"serial_number"`   // SN / Token / Kode Voucher (Penting untuk produk game/pulsa)
		Price          float64           `json:"price"`           //nominal tagihan
		AdminFee       float64           `json:"AdminFee"`        //nominal admin dari provider
		MerchantFee    float64           `json:"MerchantFee"`     //nominal fee dari provider ke kita
		GrandTotal     float64           `json:"grand_total"`     //total pembayaran yang perlu dibayar konsumen
		LastBalance    float64           `json:"last_balance"`    // Sisa saldo deposit di provider (jika ada)
		ProviderDetail ProviderFeedback  `json:"provider_detail"` //deskripsi status dari provider
		BillDesc       string            `json:"bill_desc"`       //kumpulan informasi terkait product format baku
		ProcessedAt    string            `json:"processed_at"`
	}

	// ProviderFeedback menampung data mentah dari provider tanpa merusak format standar
	ProviderFeedback struct {
		Code        string `json:"code"`        // RC Asli Provider (misal: "00", "141", "39")
		Message     string `json:"message"`     // Message Murni Provider
		Description string `json:"description"` // Description Murni Provider
	}
)
type (
	InquiryResult struct {
		CustomerID     string           `json:"customer_id"`
		RefID          string           `json:"ref_id"`          // Reference ID dari Main Service
		ProviderRefID  string           `json:"provider_ref_id"` // Transaction ID / Invoice ID dari Provider
		Price          float64          `json:"price"`
		AdminFee       float64          `json:"AdminFee"`
		MerchantFee    float64          `json:"MerchantFee"`
		GrandTotal     float64          `json:"grand_total"`     //total pembayaran yang perlu dibayar konsumen
		ProviderDetail ProviderFeedback `json:"provider_detail"` //deskripsi status dari provider
		BillDesc       string           `json:"bill_desc"`       //kumpulan informasi terkait product format baku
		ProcessedAt    string           `json:"processed_at"`
	}
)

type (
	RequestInquiry struct {
		RefID       string      `json:"ref_id"` // Reference ID dari Main Service
		CustomerID  string      `json:"customer_id"`
		DataProduct DataProduct `json:"data_product"`
	}
	RequestPayment struct {
		RefID         string `json:"ref_id"`          // Reference ID dari Main Service
		ProviderRefID string `json:"provider_ref_id"` // Transaction ID / Invoice ID dari Provider
	}
	DataProduct struct {
		ProviderID        int64  `json:"provider_id"`
		ProductCategoryID int64  `json:"product_category_id"`
		ProductTypeID     int64  `json:"product_type_id"`
		ProductCode       string `json:"product_code"`
	}
)
