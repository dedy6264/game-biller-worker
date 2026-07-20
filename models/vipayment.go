package models

type (
	ReqPrepaidPaymentVIP struct {
		Key     string `json:"key"`
		Sign    string `json:"sign"`
		Type    string `json:"type"`
		Service string `json:"service"`
		DataNo  string `json:"data_no"`
		TrxId   string `json:"trxid"`
		Limit   string `json:"limit"`
	}
	RespPrepaidPaymentVIP struct {
		Result bool `json:"result"`
		Data   struct {
			Trxid   string `json:"trxid"`
			Data    string `json:"data"`
			Code    string `json:"code"`
			Service string `json:"service"`
			Status  string `json:"status"`
			Note    string `json:"note"`
			Balance int    `json:"balance"`
			Price   int    `json:"price"`
		} `json:"data"`
		Message string `json:"message"`
	}
)
