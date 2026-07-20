package models

type (
	ReqInquiryPostpaidIAK struct {
		Commands string `json:"commands"`
		Hp       string `json:"hp"`
		Code     string `json:"code"`
		RefId    string `json:"ref_id"`
		Username string `json:"username"`
		Sign     string `json:"sign"`
		Month    string `json:"month"`
	}
	ReqPaymentPostpaidIAK struct {
		Commands string `json:"commands"`
		Username string `json:"username"`
		TrID     string `json:"tr_id"`
		Sign     string `json:"sign"`
	}
	ReqPaymentPrepaidIAK struct { //format payment V2
		CustomerId  string `json:"customer_id"`
		ProductCode string `json:"product_code"`
		RefId       string `json:"ref_id"`
		Username    string `json:"username"`
		Sign        string `json:"sign"`
	}
	ReqInquiryPlnTokenIAK struct {
		Username   string `json:"username"`
		CustomerID string `json:"customer_id"`
		Sign       string `json:"sign"`
	}
	ReqInquiryPostpaidIak struct {
		Commands string `json:"commands"`
		Hp       string `json:"hp"`
		Code     string `json:"code"`
		RefId    string `json:"ref_id"`
		Username string `json:"username"`
		Sign     string `json:"sign"`
		Month    string `json:"month"`
	}
)
type (
	RespInquiryBpjsIAK struct {
		Data struct {
			TrID         int    `json:"tr_id"`
			Code         string `json:"code"`
			Hp           string `json:"hp"`
			TrName       string `json:"tr_name"`
			Period       string `json:"period"`
			Nominal      int    `json:"nominal"`
			Admin        int    `json:"admin"`
			RefID        string `json:"ref_id"`
			ResponseCode string `json:"response_code"`
			Message      string `json:"message"`
			Price        int    `json:"price"`
			SellingPrice int    `json:"selling_price"`
			Desc         struct {
				KodeCabang     string `json:"kode_cabang"`
				NamaCabang     string `json:"nama_cabang"`
				SisaPembayaran string `json:"sisa_pembayaran"`
				JumlahPeserta  string `json:"jumlah_peserta"`
			} `json:"desc"`
		} `json:"data"`
		Meta []interface{} `json:"meta"`
	}
	RespInquiryPlnTokenIAK struct {
		Data struct {
			Status       any    `json:"status"`
			CustomerID   string `json:"customer_id"`
			MeterNo      string `json:"meter_no"`
			SubscriberID string `json:"subscriber_id"`
			Name         string `json:"name"`
			SegmentPower string `json:"segment_power"`
			Message      string `json:"message"`
			Rc           string `json:"rc"`
		} `json:"data"`
		Meta []interface{} `json:"meta"`
	}
	RespInquiryPlnIAK struct {
		Data struct {
			TrID         int    `json:"tr_id"`
			Code         string `json:"code"`
			Hp           string `json:"hp"`
			TrName       string `json:"tr_name"`
			Period       string `json:"period"`
			Nominal      int    `json:"nominal"`
			Admin        int    `json:"admin"`
			RefID        string `json:"ref_id"`
			ResponseCode string `json:"response_code"`
			Message      string `json:"message"`
			Price        int    `json:"price"`
			SellingPrice int    `json:"selling_price"`
			Desc         struct {
				Tarif         string `json:"tarif"`
				Daya          int    `json:"daya"`
				LembarTagihan string `json:"lembar_tagihan"`
				Tagihan       struct {
					Detail []struct {
						Periode      string `json:"periode"`
						NilaiTagihan string `json:"nilai_tagihan"`
						Admin        string `json:"admin"`
						Denda        string `json:"denda"`
						Total        int    `json:"total"`
					} `json:"detail"`
				} `json:"tagihan"`
			} `json:"desc"`
		} `json:"data"`
		Meta []interface{} `json:"meta"`
	}
	RespInquiryFinanceIAK struct {
		Data struct {
			TrID         int    `json:"tr_id"`
			Code         string `json:"code"`
			Hp           string `json:"hp"`
			TrName       string `json:"tr_name"`
			Period       string `json:"period"`
			Nominal      int    `json:"nominal"`
			Admin        int    `json:"admin"`
			RefID        string `json:"ref_id"`
			ResponseCode string `json:"response_code"`
			Message      string `json:"message"`
			Price        int    `json:"price"`
			SellingPrice int    `json:"selling_price"`
			Desc         struct {
				MiscFee         int    `json:"misc_fee"`
				ItemName        string `json:"item_name"`
				NoRankga        string `json:"no_rankga"`
				NoPol           string `json:"no_pol"`
				Tenor           string `json:"tenor"`
				Installment     int    `json:"installment"`
				PenaltyBill     int    `json:"penalty_bill"`
				MaxPayment      int    `json:"max_payment"`
				LastPaidDueDate string `json:"last_paid_due_date"`
				IDRef           string `json:"id_ref"`
			} `json:"desc"`
		} `json:"data"`
		Meta []interface{} `json:"meta"`
	}
	RespPaymentFinanceIAK struct {
		Data struct {
			TrID         int    `json:"tr_id"`
			Code         string `json:"code"`
			Datetime     string `json:"datetime"`
			Hp           string `json:"hp"`
			TrName       string `json:"tr_name"`
			Period       string `json:"period"`
			Nominal      int    `json:"nominal"`
			Admin        int    `json:"admin"`
			ResponseCode string `json:"response_code"`
			Message      string `json:"message"`
			Price        int    `json:"price"`
			SellingPrice int    `json:"selling_price"`
			Balance      int    `json:"balance"`
			Noref        string `json:"noref"`
			RefID        string `json:"ref_id"`
			Desc         struct {
				MiscFee         int    `json:"misc_fee"`
				ItemName        string `json:"item_name"`
				NoRankga        string `json:"no_rankga"`
				NoPol           string `json:"no_pol"`
				Tenor           string `json:"tenor"`
				Installment     int    `json:"installment"`
				PenaltyBill     int    `json:"penalty_bill"`
				MaxPayment      int    `json:"max_payment"`
				LastPaidDueDate string `json:"last_paid_due_date"`
				IDRef           string `json:"id_ref"`
			} `json:"desc"`
		} `json:"data"`
		Meta []interface{} `json:"meta"`
	}
	RespPaymentPlnIAK struct {
		Data struct {
			TrID         int    `json:"tr_id"`
			Code         string `json:"code"`
			Datetime     string `json:"datetime"`
			Hp           string `json:"hp"`
			TrName       string `json:"tr_name"`
			Period       string `json:"period"`
			Nominal      int    `json:"nominal"`
			Admin        int    `json:"admin"`
			ResponseCode string `json:"response_code"`
			Message      string `json:"message"`
			Price        int    `json:"price"`
			SellingPrice int    `json:"selling_price"`
			Balance      int    `json:"balance"`
			Noref        string `json:"noref"`
			RefID        string `json:"ref_id"`
			Desc         struct {
				Tarif             string `json:"tarif"`
				Daya              int    `json:"daya"`
				LembarTagihan     string `json:"lembar_tagihan"`
				LembarTagihanSisa int    `json:"lembar_tagihan_sisa"`
				Tagihan           struct {
					Detail []struct {
						MeterAwal    string `json:"meter_awal"`
						MeterAkhir   string `json:"meter_akhir"`
						Periode      string `json:"periode"`
						NilaiTagihan string `json:"nilai_tagihan"`
						Admin        string `json:"admin"`
						Denda        string `json:"denda"`
						Total        int    `json:"total"`
					} `json:"detail"`
				} `json:"tagihan"`
			} `json:"desc"`
		} `json:"data"`
		Meta []interface{} `json:"meta"`
	}
	RespPaymentBPJSIAK struct {
		Data struct {
			TrID         int    `json:"tr_id"`
			Code         string `json:"code"`
			Datetime     string `json:"datetime"`
			Hp           string `json:"hp"`
			TrName       string `json:"tr_name"`
			Period       string `json:"period"`
			Nominal      int    `json:"nominal"`
			Admin        int    `json:"admin"`
			ResponseCode string `json:"response_code"`
			Message      string `json:"message"`
			Price        int    `json:"price"`
			SellingPrice int    `json:"selling_price"`
			Balance      int    `json:"balance"`
			Noref        string `json:"noref"`
			RefID        string `json:"ref_id"`
			Desc         struct {
				KodeCabang     string `json:"kode_cabang"`
				NamaCabang     string `json:"nama_cabang"`
				SisaPembayaran string `json:"sisa_pembayaran"`
				JumlahPeserta  string `json:"jumlah_peserta"`
			} `json:"desc"`
		} `json:"data"`
		Meta []interface{} `json:"meta"`
	}
	RespPaymentPrepaidIAK struct {
		Data struct {
			TrID         int    `json:"tr_id"`
			Message      string `json:"message"`
			Price        int    `json:"price"`
			Balance      int    `json:"balance"`
			RefID        string `json:"ref_id"`
			Status       int    `json:"status"`
			ProductCode  string `json:"product_code"`
			CustomerID   string `json:"customer_id"`
			SubscriberID string `json:"subscriber_id"` // untuk fallback cek PLN prepaid
			Rc           string `json:"rc"`
			Sn           string `json:"sn"`
		} `json:"data"`
		Meta []interface{} `json:"meta"`
	}
	RespPaymentPostpaidIAK struct {
		Data struct {
			TrID         int    `json:"tr_id"`
			Code         string `json:"code"`
			Hp           string `json:"hp"`
			TrName       string `json:"tr_name"`
			Period       string `json:"period"`
			Nominal      int    `json:"nominal"`
			Admin        int    `json:"admin"`
			RefID        string `json:"ref_id"`
			ResponseCode string `json:"response_code"`
			Message      string `json:"message"`
			Price        int    `json:"price"`
			SellingPrice int    `json:"selling_price"`
			Rc           string `json:"rc"` // field normalisasi
		} `json:"data"`
		Meta []interface{} `json:"meta"`
	}
)
type (
	RespWorkerIakUndefined struct {
		ResponseCode string `json:"response_code"`
		Message      string `json:"message"`
	}
	RespWorkerIakUndefinedI struct {
		Data struct {
			ResponseCode string `json:"response_code"`
			Message      string `json:"message"`
		} `json:"data"`
	}
	RespWorkerIakUndefinedII struct {
		Data struct {
			Rc      string `json:"rc"`
			Message string `json:"message"`
			Status  int    `json:"status"`
		} `json:"data"`
	}
)
