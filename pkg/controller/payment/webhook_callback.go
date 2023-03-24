package payment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/services/payment"
	"github.com/vesicash/payment-ms/utility"
)

func (base *Controller) RaveWebhook(c *gin.Context) {
	var (
		req models.RaveWebhookRequest
	)

	err := c.ShouldBind(&req)
	if err != nil {
		base.Logger.Error("rave webhhook log error", "Failed to parse request body", err.Error())
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	code, err := payment.RaveWebhookService(c, base.ExtReq, base.Db, req)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "ok", nil)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) MonnifyWebhook(c *gin.Context) {
	var (
		req models.MonnifyWebhookRequest
	)

	err := c.ShouldBind(&req)
	if err != nil {
		base.Logger.Error("monnify webhhook log error", "Failed to parse request body", err.Error())
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	code, err := payment.MonnifyWebhookService(c, base.ExtReq, base.Db, req)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "ok", nil)
	c.JSON(http.StatusOK, rd)

}
func (base *Controller) MonnifyDisbursementCallback(c *gin.Context) {
	var (
		req models.MonnifyWebhookRequest
	)

	err := c.ShouldBind(&req)
	if err != nil {
		base.Logger.Error("monnify callback log error", "Failed to parse request body", err.Error())
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	code, err := payment.MonnifyDisbursementCallbackService(c, base.ExtReq, base.Db, req)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "ok", nil)
	c.JSON(http.StatusOK, rd)

}
