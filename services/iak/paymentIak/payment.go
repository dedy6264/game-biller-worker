package paymentiak

import (
	"encoding/json"
	"game-biller-worker/constans"
	"game-biller-worker/helpers"
	"game-biller-worker/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Payment(c echo.Context) error {
	var (
		svc         = "IAK Payment"
		request     models.RequestPayment
		otherCustId models.OtherCustomerID
	)

	// Bind request body ke models.RequestPayment
	if err := c.Bind(&request); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusBadRequest, models.PaymentResult{})
	}

	var (
		result models.PaymentResult
		err    error
	)
	isPLN := request.DataProduct.ProductCategoryID == constans.PRODUCT_CATEGORY_PLN
	isPrepaid := request.DataProduct.ProductTypeID == constans.PRODUCT_TYPE_PREPAID
	// isPostpaid := request.DataProduct.ProductTypeID == constans.PRODUCT_TYPE_POSTPAID
	if isPrepaid {
		if isPLN {
			result, err = PlnToken(request)
		} else {
			switch request.DataProduct.ProductReferenceID {
			case 8: //ML
				_ = json.Unmarshal([]byte(request.OtherCustomerID), &otherCustId)
				request.CustomerID = request.CustomerID + "|" + otherCustId.ZoneID
			case 9: //GENSHIN IMPACT
				_ = json.Unmarshal([]byte(request.OtherCustomerID), &otherCustId)
				request.CustomerID = request.CustomerID + "|" + otherCustId.ServerID
			case 12: //ragnarok
				_ = json.Unmarshal([]byte(request.OtherCustomerID), &otherCustId)
				request.CustomerID = request.CustomerID + "|" + otherCustId.ServerID
			case 13, 10, 14, 15: //Point Blank
				request.CustomerID = request.CustomerID
				// case 10: //Free Fire
				// 	_ = json.Unmarshal([]byte(request.OtherCustomerID), &otherCustId)
				// 	request.CustomerID = request.CustomerID + "|" + otherCustId.ServerID
			}

			result, err = Prepaid(request)
		}
	} else {
		if isPLN {
			result, err = PlnPostpaid(request)
		} else {

		}
	}

	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Sub-function error")
		return c.JSON(http.StatusInternalServerError, models.PaymentResult{})
	}

	return c.JSON(http.StatusOK, result)
}
