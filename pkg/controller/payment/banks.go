package payment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/payment-ms/services/payment"
	"github.com/vesicash/payment-ms/utility"
)

func (base *Controller) ListBanks(c *gin.Context) {
	var (
		req struct {
			CountryCode string `json:"country_code" validate:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	banks, code, err := payment.ListBanksService(base.ExtReq, base.Db, req.CountryCode)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Banks fetched", banks)
	c.JSON(http.StatusOK, rd)

}
func (base *Controller) VerifyBankAccount(c *gin.Context) {
	var (
		req struct {
			BankCode      string `json:"bank_code" validate:"required"`
			AccountNumber string `json:"account_number" validate:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	banks, code, err := payment.VerifyBankAccountService(base.ExtReq, base.Db, req.BankCode, req.AccountNumber)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Account number is valid", banks)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) ConvertCurrency(c *gin.Context) {
	var (
		req struct {
			Amount float64 `json:"amount" validate:"required"`
			From   string  `json:"from" validate:"required"`
			To     string  `json:"to" validate:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	data, code, err := payment.ConvertCurrencyService(base.ExtReq, base.Db, req.Amount, req.From, req.To)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Conversion Successful", data)
	c.JSON(http.StatusOK, rd)

}
