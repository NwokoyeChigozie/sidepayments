package payment

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/services/payment"
	"github.com/vesicash/payment-ms/utility"
)

func (base *Controller) GetCustomerPayments(c *gin.Context) {
	var (
		businessID = c.Param("business_id")
	)

	businessIDInt, err := strconv.Atoi(businessID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "incorrect business id format", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	total, code, err := payment.GetCustomerPaymentsService(base.ExtReq, base.Db, businessIDInt)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Data Retrieved", gin.H{"total": total})
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) GetCustomerCard(c *gin.Context) {

	user := models.MyIdentity
	if user == nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "error retrieving authenticated user", fmt.Errorf("error retrieving authenticated user"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	card, code, err := payment.GetCustomerCardService(base.ExtReq, base.Db, user.AccountID)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Data Retrieved", card)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) GetCustomerCardsByBusinessID(c *gin.Context) {
	var (
		businessID = c.Param("business_id")
	)

	businessIDInt, err := strconv.Atoi(businessID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "incorrect business id format", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	card, code, err := payment.GetCustomerCardsByBusinessIDService(base.ExtReq, base.Db, businessIDInt)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Data Retrieved", card)
	c.JSON(http.StatusOK, rd)

}
