package payment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/payment-ms/services/payment"
	"github.com/vesicash/payment-ms/utility"
)

func (base *Controller) GetPaymentInvoice(c *gin.Context) {
	var (
		paymentID = c.Param("payment_id")
	)

	template, data, code, err := payment.GetPaymentInvoiceService(c, base.ExtReq, base.Db, paymentID)
	if err != nil {
		c.String(code, err.Error())
		return
	}

	err = template.Execute(c.Writer, data)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (base *Controller) RenderPayStatus(c *gin.Context) {
	var (
		status  = c.Param("status")
		website = c.Param("website")
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
