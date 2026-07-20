package models

type (
	ReqHeader struct {
		Header []Header
	}

	Header struct {
		Key      string
		Val      string
		IsUpCase bool
	}
)
