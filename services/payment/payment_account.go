package payment

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/config"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/payment-ms/utility"
)

func PaymentAccountMonnifyListService(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, req models.PaymentAccountMonnifyListRequest) (interface{}, int, error) {
	var (
		data               interface{}
		generatedReference = ""
		configData = config.GetConfig()
	)

	user, err := GetUserWithAccountID(extReq, req.AccountID)
	if err != nil {
		return data, http.StatusInternalServerError, err
	}

	businessProfile, err := GetBusinessProfileByAccountID(extReq, extReq.Logger, req.AccountID)
	if err != nil {
		return data, http.StatusInternalServerError, fmt.Errorf("business profile not found for user: %v", err)
	}

	accountName := businessProfile.BusinessName
	if accountName == "" {
		accountName = user.Firstname
		if accountName == "" {
			accountName = strconv.Itoa(int(user.AccountID))
		}
	}
	accountEmail := user.EmailAddress

	if req.TransactionID != "" {
		transaction, err := ListTransactionsByID(extReq, req.TransactionID)
		if err != nil {
			return data, http.StatusInternalServerError, err
		}
		paymentAccounts := models.PaymentAccount{BusinessID: strconv.Itoa(int(user.AccountID)), TransactionID: req.TransactionID}
		code, err := paymentAccounts.GetPaymentAccountByBusinessIDAndTransactionID(db.Payment)
		if err != nil {
			if code == http.StatusInternalServerError {
				return data, code, err
			}
			if req.GeneratedReference != "" {
				generatedReference = req.GeneratedReference
			} else {
				uuID, _ := uuid.NewV4()
				generatedReference = "VESICASH_VA_" + uuID.String()
			}
			currencyCode := "NGN"

			paymentInfo := models.PaymentInfo{Reference: generatedReference}
			code, err := paymentInfo.GetPaymentInfoByReference(db.Payment)
			if err != nil {
				if code == http.StatusInternalServerError {
					return data, code, err
				}

				payment := models.Payment{
					PaymentID:    utility.RandomString(10),
					TotalAmount:  transaction.Amount,
					EscrowCharge: transaction.EscrowCharge,
					IsPaid:       false,
					AccountID:    int64(req.AccountID),
					BusinessID:   int64(businessProfile.AccountID),
					Currency:     currencyCode,
				}

				err = payment.CreatePayment(db.Payment)
				if err != nil {
					return data, http.StatusInternalServerError, err
				}

				paymentInfo = models.PaymentInfo{
					PaymentID: payment.PaymentID,
					Reference: generatedReference,
					Status:    "pending",
					Gateway:   req.Gateway,
				}
				err = paymentInfo.CreatePaymentInfo(db.Payment)
				if err != nil {
					return data, http.StatusInternalServerError, err
				}
			}

			paymentAccounts := models.PaymentAccount{
				PaymentAccountID: generatedReference,
				TransactionID: req.TransactionID,
				PaymentID: paymentInfo.PaymentID,
			}

			if req.Gateway == "rave" {
				paymentAccounts.AccountNumber = 
				// $payment_accounts->account_number = config("payment.rave.merchant_id");
//                                 $payment_accounts->account_name = config("payment.rave.account_name");
//                                 $payment_accounts->bank_code = "flutterwave";
//                                 $payment_accounts->bank_name = "flutterwave";
//                                 $payment_accounts->status = "ACTIVE";
			}else {

			}




		}
	}

}

// $generatedReference = "";
//                 if (isset($request->transaction_id)){
//                     $transactionId = $request->transaction_id;
//                     $paymentAccountCheck = PaymentAccounts::where('business_id', $account_id)
//                         ->where("transaction_id",$transactionId)->count();

//                     if ($paymentAccountCheck  == 0){
//                         if (isset($request->generated_reference)) {
//                             $generatedReference = $request->generated_reference;
//                         }else{
//                             $generatedReference = "VESICASH_VA_" . $this->generateUUID();
//                         }

//                         $accountName = $account_name;
//                         $currencyCode = "NGN";
//                         $customerEmail = $account_email;

//                         // create payment info if not exist
//                     $paymentInfo = null;
//                         // Store payment info
//                         $payment_info_check = PaymentInfo::where('reference', $generatedReference)->count();
//                         if ($payment_info_check == 0) {
//                             $payment_info = new PaymentInfo();
//                             $payment_info->payment_id = str_random(10);
//                             $payment_info->reference  = $generatedReference;
//                             $payment_info->status = 'pending';
//                             $payment_info->gateway = $request->gateway ?? 'monnify';
//                             $payment_info->save();
//                             $paymentInfo = $payment_info;
//                         }

//                         // account data
//                         $data = [
//                             'reference'     => $generatedReference,
//                             'accountName'   => $accountName,
//                             'currencyCode'  => $currencyCode,
//                             'customerEmail' => $customerEmail
//                         ];

//                         // reserve an account

//                         $payment_accounts = new PaymentAccounts();
//                         $payment_accounts->payment_account_id = $generatedReference;
//                         $payment_accounts->transaction_id = $transactionId;
//                         $payment_accounts->payment_id = $paymentInfo->payment_id;

//                         if (isset($request->gateway)){
//                             $gateway = $request->gateway;
//                             if ($gateway == "rave"){
//                                 $payment_accounts->account_number = config("payment.rave.merchant_id");
//                                 $payment_accounts->account_name = config("payment.rave.account_name");
//                                 $payment_accounts->bank_code = "flutterwave";
//                                 $payment_accounts->bank_name = "flutterwave";
//                                 $payment_accounts->status = "ACTIVE";
//                             }
//                         }else{
//                             // call on Monnify endpoint.
//                             $monnify = new Monnify();
//                             $reserveAccount = $monnify->reserveAccount($data);
//                             $accountDetails = $reserveAccount;

//                             $payment_accounts->account_number = $accountDetails->accountNumber;
//                             $payment_accounts->account_name = $accountDetails->accountName;
//                             $payment_accounts->bank_code = $accountDetails->bankCode;
//                             $payment_accounts->bank_name = $accountDetails->bankName;
//                             $payment_accounts->reservation_reference = $accountDetails->reservationReference;
//                             $payment_accounts->status = $accountDetails->status;
//                         }

//                         $payment_accounts->is_used = true;
//                         $payment_accounts->expires_after = strtotime("+30 days", time());
//                         $payment_accounts->business_id = $account_id;
//                         $payment_accounts->save();

//                         $slackPaymentChannelEnv = env('SLACK_PAYMENT_CHANNEL', 'payments-test');
//                         try {
//                             // Notify on Slack
//                             $this->slackNotify($slackPaymentChannelEnv, '
//                     Virtual Bank Account Generated
//                     ```
//                     Environment: ' . env('APP_NAME') . '
//                     Account Number: ' . $payment_accounts->account_number . '
//                     Account Name: ' . $payment_accounts->account_name . '
//                     Bank: ' . $payment_accounts->bank_name . '
//                     Status: SUCCESSFUL
//                         ```');
//                         } catch (\Exception $e) {
//                             return $this->response('error', 'Slack notification error', $e->getMessage(), 500);
//                         }
//                    here }else{
//                         // return a reserved account
//                         $payment_accounts = PaymentAccounts::where('business_id', $account_id)
//                             ->where("transaction_id",$transactionId)
//                             ->first();

//                         $transactionId =  $request->transaction_id;
//                         $payment_accounts->update([
//                             "transaction_id" => $transactionId
//                         ]);
//                     }
//                 }
