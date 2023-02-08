package payment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/payment-ms/services/payment"
	"github.com/vesicash/payment-ms/utility"
)

func (base *Controller) ListPayments(c *gin.Context) {
	var (
		transactionID = c.Param("transaction_id")
	)

	payments, code, err := payment.LisPaymentsService(base.ExtReq, base.Db, transactionID)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", payments)
	c.JSON(http.StatusOK, rd)

}
