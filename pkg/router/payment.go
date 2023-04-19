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

	paymentUrl := r.Group(fmt.Sprintf("%v", ApiVersion))
	{
		paymentUrl.POST("/banks", payment.ListBanks)
		paymentUrl.POST("/currency/converter", payment.ConvertCurrency)

		paymentUrl.POST("/webhook/rave", payment.RaveWebhook)
		paymentUrl.POST("/webhook/monnify", payment.MonnifyWebhook)
		paymentUrl.POST("/disbursement/callback", payment.MonnifyDisbursementCallback)
		paymentUrl.GET("/disbursement/callback", payment.MonnifyDisbursementCallback)

		paymentUrl.GET("/payment/invoice/:payment_id", payment.GetPaymentInvoice)
		paymentUrl.GET("/pay/successful", payment.RenderPaySuccessful)
		paymentUrl.GET("/pay/failed", payment.RenderPayFailed)
		paymentUrl.GET("/pay/status", payment.GetStatus)
		paymentUrl.POST("/pay/new-status", payment.GetPaymentStatus)
	}

	paymentAuthUrl := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db, extReq, middleware.AuthType))
	{
		paymentAuthUrl.POST("/create", payment.CreatePayment)
		paymentAuthUrl.PATCH("/edit", payment.EditPayment)
		paymentAuthUrl.DELETE("/delete/:id", payment.DeletePayment)
		paymentAuthUrl.GET("/customers/card", payment.GetCustomerCard)
		paymentAuthUrl.GET("/customers/cards/:business_id", payment.GetCustomerCardsByBusinessID)
	}

	paymentApiUrl := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db, extReq, middleware.ApiType))
	{
		paymentApiUrl.POST("/create/headless", payment.CreatePaymentHeadless)
		paymentApiUrl.GET("/listByPaymentId/:payment_id", payment.GetPaymentByID)
		paymentApiUrl.GET("/listByTransactionId/:transaction_id", payment.ListPaymentByTransactionID)
		paymentApiUrl.GET("/list-payments/:transaction_id", payment.ListPaymentRecords)
		paymentApiUrl.GET("/list/wallet_funding/:account_id", payment.ListPaymentsByAccountID)
		paymentApiUrl.GET("/list/wallet_withdrawals/:account_id", payment.ListWithdrawalsByAccountID)
		paymentApiUrl.POST("/verify", payment.VerifyTransactionPayment)
		paymentApiUrl.GET("/customers/payments/:business_id", payment.GetCustomerPayments)
		// /pay
		paymentApiUrl.POST("/pay", payment.InitiatePayment)
		paymentApiUrl.POST("/pay/headless", payment.InitiatePaymentHeadless)
		paymentApiUrl.POST("/pay/tokenized", payment.ChargeCardInit)
		paymentApiUrl.POST("/pay/tokenized/headless", payment.ChargeCardHeadlessInit)
		paymentApiUrl.DELETE("/pay/tokenized/delete", payment.DeleteStoredCard)

		paymentApiUrl.POST("/payment_account/list", payment.PaymentAccountMonnifyList)
		paymentApiUrl.POST("/payment_account/verify", payment.PaymentAccountMonnifyVerify)

		paymentApiUrl.GET("/disbursement/user/:account_id", payment.ListDisbursementByAccountID)
		paymentApiUrl.POST("/disbursement/wallet/wallet-transfer", payment.WalletTransfer)
		paymentApiUrl.POST("/disbursement/wallet/withdraw", payment.ManualDebit)
		paymentApiUrl.POST("/disbursement/process/refund", payment.ManualRefund)

	}

	paymentAppUrl := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db, extReq, middleware.AppType))
	{
		paymentAppUrl.POST("/wallet/debit", payment.DebitWallet)
		paymentAppUrl.POST("/wallet/credit", payment.CreditWallet)
	}

	paymentjobsUrl := r.Group(fmt.Sprintf("%v/jobs", ApiVersion))
	{
		paymentjobsUrl.POST("/start", payment.StartCronJob)
		paymentjobsUrl.POST("/start-bulk", payment.StartCronJobsBulk)
		paymentjobsUrl.POST("/stop", payment.StopCronJob)
		paymentjobsUrl.PATCH("/update_interval", payment.UpdateCronJobInterval)
	}
	return r
}
