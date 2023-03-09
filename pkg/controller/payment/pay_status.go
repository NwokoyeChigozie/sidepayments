package payment

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/payment-ms/services/payment"
	"github.com/vesicash/payment-ms/utility"
)

func (base *Controller) GetStatus(c *gin.Context) {
	var (
		reference        = c.Query("reference")
		headlessString   = c.Query("headless")
		successPage      = c.Query("success_page")
		failurePage      = c.Query("failure_page")
		fundWalletString = c.Query("fund_wallet")
	)

	req := models.GetStatusRequest{Reference: reference, SuccessPage: successPage, FailurePage: failurePage}
	if headlessString != "" && strings.ToLower(headlessString) == "true" {
		req.Headless = true
	} else {
		req.Headless = false
	}
	if fundWalletString != "" && strings.ToLower(fundWalletString) == "true" {
		req.FundWallet = true
	} else {
		req.FundWallet = false
	}

	err := base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	vr := postgresql.ValidateRequestM{Logger: base.Logger, Test: base.ExtReq.Test}
	err = vr.ValidateRequest(req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	msg, code, err := payment.GetStatusService(c, base.ExtReq, base.Db, req)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, msg, nil)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) GetPaymentStatus(c *gin.Context) {
	var (
		req models.GetPaymentStatusRequest
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

	vr := postgresql.ValidateRequestM{Logger: base.Logger, Test: base.ExtReq.Test}
	err = vr.ValidateRequest(req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	uri, msg, code, err := payment.GetPaymentStatusService(c, base.ExtReq, base.Db, req)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, msg, uri)
	c.JSON(http.StatusOK, rd)

}
