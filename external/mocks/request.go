package mocks

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/mocks/appruve_mocks"
	"github.com/vesicash/payment-ms/external/mocks/auth_mocks"
	"github.com/vesicash/payment-ms/external/mocks/ip_api_mocks"
	"github.com/vesicash/payment-ms/external/mocks/ipstack_mocks"
	"github.com/vesicash/payment-ms/external/mocks/monnify_mocks"
	"github.com/vesicash/payment-ms/external/mocks/notification_mocks"
	"github.com/vesicash/payment-ms/external/mocks/rave_mocks"
	"github.com/vesicash/payment-ms/external/mocks/transactions_mocks"
	"github.com/vesicash/payment-ms/external/mocks/upload_mocks"
	"github.com/vesicash/payment-ms/utility"
)

type ExternalRequest struct {
	Logger     *utility.Logger
	Test       bool
	RequestObj RequestObj
}

type RequestObj struct {
	Name         string
	Path         string
	Method       string
	Headers      map[string]string
	SuccessCode  int
	RequestData  interface{}
	DecodeMethod string
	Logger       *utility.Logger
}

var (
	JsonDecodeMethod    string = "json"
	PhpSerializerMethod string = "phpserializer"
)

func (er ExternalRequest) SendExternalRequest(name string, data interface{}) (interface{}, error) {
	switch name {
	case "get_user":
		return auth_mocks.GetUser(er.Logger, data)
	case "get_user_credential":
		return auth_mocks.GetUserCredential(er.Logger, data)
	case "create_user_credential":
		return auth_mocks.CreateUserCredential(er.Logger, data)
	case "update_user_credential":
		return auth_mocks.UpdateUserCredential(er.Logger, data)
	case "get_user_profile":
		return auth_mocks.GetUserProfile(er.Logger, data)
	case "get_business_profile":
		return auth_mocks.GetBusinessProfile(er.Logger, data)
	case "get_country":
		return auth_mocks.GetCountry(er.Logger, data)
	case "get_bank_details":
		return auth_mocks.GetBankDetails(er.Logger, data)
	case "get_access_token":
		return auth_mocks.GetAccessToken(er.Logger)
	case "validate_on_auth":
		return auth_mocks.ValidateOnAuth(er.Logger, data)
	case "validate_authorization":
		return auth_mocks.ValidateAuthorization(er.Logger, data)
	case "send_verification_email":
		return notification_mocks.SendVerificationEmail(er.Logger, data)
	case "send_welcome_email":
		return notification_mocks.SendWelcomeEmail(er.Logger, data)
	case "send_email_verified_notification":
		return notification_mocks.SendEmailVerifiedNotification(er.Logger, data)
	case "send_sms_to_phone":
		return notification_mocks.SendSendSMSToPhone(er.Logger, data)
	case "monnify_login":
		return monnify_mocks.MonnifyLogin(er.Logger, data)
	case "monnify_match_bvn_details":
		return monnify_mocks.MonnifyMatchBvnDetails(er.Logger, data)
	case "appruve_verify_id":
		return appruve_mocks.AppruveVerifyID(er.Logger, data)
	case "verification_failed_notification":
		return notification_mocks.VerificationFailedNotification(er.Logger, data)
	case "verification_successful_notification":
		return notification_mocks.VerificationSuccessfulNotification(er.Logger, data)
	case "rave_resolve_bank_account":
		return rave_mocks.RaveResolveBankAccount(er.Logger, data)
	case "ipstack_resolve_ip":
		return ipstack_mocks.IpstackResolveIp(er.Logger, data)
	case "get_authorize":
		return auth_mocks.GetAuthorize(er.Logger, data)
	case "create_authorize":
		return auth_mocks.CreateAuthorize(er.Logger, data)
	case "update_authorize":
		return auth_mocks.UpdateAuthorize(er.Logger, data)
	case "send_authorized_notification":
		return notification_mocks.SendAuthorizedNotification(er.Logger, data)
	case "send_authorization_notification":
		return notification_mocks.SendAuthorizationNotification(er.Logger, data)
	case "set_user_authorization_required_status":
		return auth_mocks.SetUserAuthorizationRequiredStatus(er.Logger, data)
	case "validate_on_transactions":
		return transactions_mocks.ValidateOnTransactions(er.Logger, data)
	case "list_transactions_by_id":
		return transactions_mocks.ListTransactionsByID(er.Logger, data)
	case "get_users_by_business_id":
		return auth_mocks.GetUsersByBusinessID(er.Logger, data)
	case "list_banks_with_rave":
		return rave_mocks.ListBanksWithRave(er.Logger, data)
	case "convert_currency_with_rave":
		return rave_mocks.ConvertCurrencyWithRave(er.Logger, data)
	case "resolve_ip":
		return ip_api_mocks.ResolveIp(er.Logger, data)
	case "get_business_charge":
		return auth_mocks.GetBusinessCharge(er.Logger, data)
	case "init_business_charge":
		return auth_mocks.InitBusinessCharge(er.Logger, data)
	case "rave_init_payment":
		return rave_mocks.RaveInitPayment(er.Logger, data)
	case "monnify_init_payment":
		return monnify_mocks.MonnifyInitPayment(er.Logger, data)
	case "get_access_token_by_key":
		return auth_mocks.GetAccessTokenByKey(er.Logger, data)
	case "get_escrow_charge":
		return transactions_mocks.GetEscrowCharge(er.Logger, data)
	case "rave_reserve_account":
		return rave_mocks.RaveReserveAccount(er.Logger, data)
	case "rave_verify_transaction_by_tx_ref":
		return rave_mocks.RaveVerifyTransactionByTxRef(er.Logger, data)
	case "monnify_verify_transaction_by_reference":
		return monnify_mocks.MonnifyVerifyTransactionByReference(er.Logger, data)
	case "create_wallet_balance":
		return auth_mocks.CreateWalletBalance(er.Logger, data)
	case "get_wallet_balance_by_account_id_and_currency":
		return auth_mocks.GetWalletBalanceByAccountIDAndCurrency(er.Logger, data)
	case "update_wallet_balance":
		return auth_mocks.UpdateWalletBalance(er.Logger, data)
	case "update_transaction_amount_paid":
		return transactions_mocks.UpdateTransactionAmountPaid(er.Logger, data)
	case "wallet_funded_notification":
		return notification_mocks.WalletfundedNotification(er.Logger, data)
	case "wallet_debit_notification":
		return notification_mocks.WalletDebitNotification(er.Logger, data)
	case "create_activity_log":
		return transactions_mocks.CreateActivityLog(er.Logger, data)
	case "payment_invoice_notification":
		notification_mocks.PaymentInvoiceNotification(er.Logger, data)
	case "transaction_update_status":
		return transactions_mocks.TransactionUpdateStatus(er.Logger, data)
	case "buyer_satisfied":
		return transactions_mocks.BuyerSatisfied(er.Logger, data)
	case "rave_charge_card":
		return rave_mocks.RaveChargeCard(er.Logger, data)
	case "monnify_reserve_account":
		return monnify_mocks.MonnifyReserveAccount(er.Logger, data)
	case "get_monnify_reserve_account_transactions":
		return monnify_mocks.GetMonnifyReserveAccountTransactions(er.Logger, data)
	case "upload_file":
		return upload_mocks.UploadFile(er.Logger, data)
	case "create_wallet_history":
		return auth_mocks.CreateWalletHistory(er.Logger, data)
	case "create_wallet_transaction":
		return auth_mocks.CreateWalletTransaction(er.Logger, data)
	case "create_exchange_transaction":
		return transactions_mocks.CreateExchangeTransaction(er.Logger, data)
	case "get_rate_by_id":
		return transactions_mocks.GetRateByID(er.Logger, data)
	case "get_bank":
		return auth_mocks.GetBank(er.Logger, data)
	case "rave_init_transfer":
		return rave_mocks.RaveInitTransfer(er.Logger, data)
	case "monnify_init_transfer":
		return monnify_mocks.MonnifyInitTransfer(er.Logger, data)
	case "transaction_paid_notification":
		return notification_mocks.TransactionPaidNotification(er.Logger, data)
	default:
		return nil, fmt.Errorf("request not found")
	}
	return nil, fmt.Errorf("request not found")
}
