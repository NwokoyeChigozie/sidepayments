package payment

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/payment-ms/services/payment"
	"github.com/vesicash/payment-ms/utility"
)

func (base *Controller) GetPaymentInvoice(c *gin.Context) {
	var (
		paymentID = c.Param("payment_id")
	)

	base.ExtReq.Logger.Error("info getting payment invoice", "payment id "+paymentID)
	template, data, code, err := payment.GetPaymentInvoiceService(c, base.ExtReq, base.Db, paymentID)
	if err != nil {
		base.ExtReq.Logger.Error("error getting payment invoice", err.Error())
		c.String(code, err.Error())
		return
	}

	base.ExtReq.Logger.Error("info getting payment invoice", fmt.Sprintf("generate template: %v, and data: %v", template, data))
	err = template.Execute(c.Writer, data)
	if err != nil {
		base.ExtReq.Logger.Error("error getting payment invoice", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	base.ExtReq.Logger.Error("info getting payment invoice", "done")
	c.Status(http.StatusOK)
}

func (base *Controller) RenderPayStatus(c *gin.Context) {
	var (
		status  = c.Param("status")
		website = c.Query("website")
	)
	if website != "" {
		url, err := utility.URLDecode(website)
		if err == nil {
			c.Redirect(http.StatusFound, url)
			return
		}
	}
	c.String(http.StatusOK, status)
}
