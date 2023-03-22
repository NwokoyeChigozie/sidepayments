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

func Payment(r *gin.Engine, ApiVersion string, validator *validator.Validate, db postgresql.Databases, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	payment := payment.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	paymentUrl := r.Group(fmt.Sprintf("%v/payment", ApiVersion))
	{
		paymentUrl.POST("/banks", payment.ListBanks)
		paymentUrl.POST("currency/converter", payment.ConvertCurrency)

		paymentUrl.POST("webhook/rave", payment.RaveWebhook)
	}

	paymentAuthUrl := r.Group(fmt.Sprintf("%v/payment", ApiVersion), middleware.Authorize(db, extReq, middleware.AuthType))
	{
		paymentAuthUrl.POST("/create", payment.CreatePayment)
		paymentAuthUrl.PATCH("/edit", payment.EditPayment)
		paymentAuthUrl.DELETE("/delete/:id", payment.DeletePayment)
		paymentAuthUrl.GET("customers/card", payment.GetCustomerCard)
		paymentAuthUrl.GET("customers/cards/:business_id", payment.GetCustomerCardsByBusinessID)

	}

	paymentApiUrl := r.Group(fmt.Sprintf("%v/payment", ApiVersion), middleware.Authorize(db, extReq, middleware.ApiType))
	{
		paymentApiUrl.POST("/create/headless", payment.CreatePaymentHeadless)
		paymentApiUrl.GET("/listByPaymentId/:payment_id", payment.GetPaymentByID)
		paymentApiUrl.GET("/listByTransactionId/:transaction_id", payment.ListPaymentByTransactionID)
		paymentApiUrl.GET("/list-payments/:transaction_id", payment.ListPaymentRecords)
		paymentApiUrl.GET("/list/wallet_funding/:account_id", payment.ListPaymentsByAccountID)
		paymentApiUrl.GET("/list/wallet_withdrawals/:account_id", payment.ListWithdrawalsByAccountID)
		paymentApiUrl.POST("/verify", payment.VerifyTransactionPayment)
		paymentApiUrl.GET("customers/payments/:business_id", payment.GetCustomerPayments)

		// /pay
		paymentAuthUrl.GET("/pay", payment.InitiatePayment)
		paymentAuthUrl.GET("/pay/headless", payment.InitiatePaymentHeadless)
		paymentAuthUrl.POST("/pay/tokenized", payment.ChargeCardInit)
		paymentAuthUrl.POST("/pay/tokenized/headless", payment.ChargeCardHeadlessInit)
		paymentAuthUrl.POST("/pay/tokenized/delete", payment.DeleteStoredCard)
		paymentAuthUrl.POST("/pay/new-status", payment.GetPaymentStatus)

		paymentAuthUrl.POST("/payment_account/list", payment.PaymentAccountMonnifyList)
		paymentAuthUrl.POST("/payment_account/verify", payment.PaymentAccountMonnifyVerify)

		paymentAuthUrl.GET("/disbursement/user/:account_id", payment.ListDisbursementByAccountID)
		paymentAuthUrl.GET("/disbursement/wallet/wallet-transfer", payment.WalletTransfer)
		paymentAuthUrl.GET("disbursement/wallet/withdraw", payment.ManualDebit)
		paymentAuthUrl.GET("disbursement/process/refund", payment.ManualRefund)

	}
	return r
}
