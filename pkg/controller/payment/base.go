package payment

import (
	"github.com/go-playground/validator/v10"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/payment-ms/utility"
)

type Controller struct {
	Db        postgresql.Databases
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}
