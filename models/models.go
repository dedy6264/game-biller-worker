package models

type (
	BillInfo struct {
		BillDesc string `json:"bill_desc"`
		Sn       string `json:"sn"`
		LembTag  int64  `json:"lemb_tag"`
	}
	BillDesc struct {
		CustomerID   string `json:"customer_id"`
		CustomerName string `json:"customer_name"`
		Periode      string `json:"periode"`
		Dependent    string `json:"dependent"`
	}
	PlnTokenBillDesc struct {
		LembTag      int64       `json:"lemb_tag"`
		CustomerID   string      `json:"customer_id"`
		CustomerName string      `json:"customer_name"`
		MeterNo      string      `json:"meter_no"`
		Tarif        string      `json:"tarif"`
		Daya         string      `json:"daya"`
		Kwh          string      `json:"kwh"`
		Details      []PlnDetail `json:"details"`
	}
	PlnDetail struct {
		Periode    string `json:"periode"`
		Tagihan    string `json:"tagihan"`
		Admin      string `json:"admin"`
		Denda      string `json:"denda"`
		MeterAwal  string `json:"meter_awal"`
		MeterAkhir string `json:"meter_akhir"`
	}
)
