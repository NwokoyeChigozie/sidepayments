package payment

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/payment-ms/services/payment"
	"github.com/vesicash/payment-ms/utility"
)

func (base *Controller) ListPaymentByTransactionID(c *gin.Context) {
	var (
		transactionID = c.Param("transaction_id")
	)

	payments, code, err := payment.ListPaymentByTransactionIDService(base.ExtReq, base.Db, transactionID)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", payments)
	c.JSON(http.StatusOK, rd)

}
func (base *Controller) ListPaymentRecords(c *gin.Context) {
	var (
		transactionID = c.Param("transaction_id")
		paginator     = postgresql.GetPagination(c)
	)

	payments, pagination, code, err := payment.ListPaymentRecordsService(base.ExtReq, base.Db, transactionID, paginator)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Payment Details Retrieved", payments, pagination)
	c.JSON(http.StatusOK, rd)

}
func (base *Controller) GetPaymentByID(c *gin.Context) {
	var (
		paymentID = c.Param("payment_id")
	)
	payment, code, err := payment.GetPaymentByIDService(base.ExtReq, base.Db, paymentID)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", payment)
	c.JSON(http.StatusOK, rd)

}
func (base *Controller) ListPaymentsByAccountID(c *gin.Context) {
	var (
		accountID = c.Param("account_id")
		paginator = postgresql.GetPagination(c)
	)

	accountIDinT, err := strconv.Atoi(accountID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "account_id provided is not integer", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	payments, pagination, code, err := payment.ListPaymentsByAccountIDService(base.ExtReq, base.Db, paginator, accountIDinT)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", payments, pagination)
	c.JSON(http.StatusOK, rd)

}
func (base *Controller) ListWithdrawalsByAccountID(c *gin.Context) {
	var (
		accountID = c.Param("account_id")
		paginator = postgresql.GetPagination(c)
	)
	accountIDinT, err := strconv.Atoi(accountID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "account_id provided is not integer", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	disbursements, pagination, code, err := payment.ListWithdrawalsByAccountIDService(base.ExtReq, base.Db, paginator, accountIDinT)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", disbursements, pagination)
	c.JSON(http.StatusOK, rd)

}
