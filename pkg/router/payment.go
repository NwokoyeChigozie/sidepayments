package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/pkg/controller/payment"
	"github.com/vesicash/payment-ms/pkg/middleware"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/payment-ms/utility"
)

func Verification(r *gin.Engine, ApiVersion string, validator *validator.Validate, db postgresql.Databases, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	payment := payment.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	paymentAuthUrl := r.Group(fmt.Sprintf("%v/payment", ApiVersion), middleware.Authorize(db, extReq, middleware.AuthType))
	{
		paymentAuthUrl.POST("/create", payment.CreatePayment)
	}
	return r
}
