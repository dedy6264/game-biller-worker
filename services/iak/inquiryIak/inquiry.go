package inquiryiak

import (
	"game-biller-worker/constans"
	"game-biller-worker/helpers"
	"game-biller-worker/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Inquiry(c echo.Context) error {
	var (
		svc     = "IAK Inquiry"
		request models.RequestInquiry
	)

	// Bind request body ke models.RequestInquiry
	if err := c.Bind(&request); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusBadRequest, models.InquiryResult{})
	}

	isPLN := request.DataProduct.ProductCategoryID == constans.PRODUCT_CATEGORY_PLN
	isPrepaid := request.DataProduct.ProductTypeID == constans.PRODUCT_TYPE_PREPAID
	isPostpaid := request.DataProduct.ProductTypeID == constans.PRODUCT_TYPE_POSTPAID

	var (
		result models.InquiryResult
		err    error
	)
	if isPLN {
		switch {
		// PLN Token (prepaid)
		case isPrepaid:
			result, err = PlnPrepaid(request)

		// PLN Pascabayar (postpaid)
		case isPostpaid:
			result, err = PlnPostpaid(request)
		default:
			return c.JSON(http.StatusBadRequest, models.InquiryResult{
				StatusCode: helpers.CodeInvalidProductSegmnt,
			})
		}
	} else {
		result = models.InquiryResult{
			StatusCode:    helpers.CodeInqSuccess,
			RefID:         request.RefID,
			ProviderRefID: "",
			DataTransaction: models.DataTransaction{
				CustomerID:   request.CustomerID,
				SerialNumber: "",
				Price:        0,
				AdminFee:     0,
				MerchantFee:  0,
				GrandTotal:   0,
				LastBalance:  0,
			},
		}
	}

	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Sub-function error")
		return c.JSON(http.StatusInternalServerError, models.InquiryResult{})
	}

	return c.JSON(http.StatusOK, result)
}
