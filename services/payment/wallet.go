package payment

import (
	"fmt"
	"strings"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/config"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/payment-ms/utility"
)

type WalletTransactionApproved string
type WalletHistoryType string

var (
	WalletTransactionYes     WalletTransactionApproved = "yes"
	WalletTransactionNo      WalletTransactionApproved = "no"
	WalletTransactionPending WalletTransactionApproved = "pending"

	WalletHistoryCredit WalletHistoryType = "credit"
	WalletHistoryDebit  WalletHistoryType = "debit"
)

func CreditWallet(extReq request.ExternalRequest, db postgresql.Databases, amount float64, currency string, businessID int, isRefund bool, creditEscrow string, transactionID string) (external_models.WalletBalance, error) {
	if currency == "" {
		currency = "NGN"
	}
	currency = strings.ToUpper(currency)

	if strings.ToLower(creditEscrow) == "yes" {
		currency = fmt.Sprintf("ESCROW_%v", currency)
		extReq.Logger.Info("credit-wallet", "creditEscrow is yes", "currency = ", currency, fmt.Sprintf("transaction with id %v", transactionID))
		if transactionID != "" && transactionID != "0" {
			_, err := UpdateTransactionAmountPaid(extReq, transactionID, amount, "+")
			if err != nil {
				extReq.Logger.Error(fmt.Sprintf("Error adding amount to transaction amount paid for transaction with id: %v, action:%v, amount:%v", transactionID, "+", amount))
			}
		}
	}

	walletBalance, err := GetWalletBalanceByAccountIdAndCurrency(extReq, businessID, currency)
	if err != nil {
		walletBalance, err = CreateWalletBalance(extReq, businessID, currency, amount)
		if err != nil {
			return walletBalance, err
		}
		extReq.Logger.Info("credit-wallet-c", "new balance:", fmt.Sprintf("%v %v", currency, amount))
	} else {
		availableBalance := walletBalance.Available + amount
		if isRefund {
			availableBalance += config.GetConfig().ONLINE_PAYMENT.DisbursementCharge
		}
		walletBalance, err = UpdateWalletBalance(extReq, walletBalance.ID, availableBalance)
		if err != nil {
			return walletBalance, err
		}

		extReq.Logger.Info("credit-wallet-u", "new balance:", fmt.Sprintf("%v %v", currency, availableBalance))
	}

	if !isRefund {
		walletEaringLog := models.WalletEarningLog{
			AccountID: businessID,
			Amount:    amount,
			Currency:  currency,
		}

		walletEaringLog.CreateWalletEarningLog(db.Payment)
		extReq.SendExternalRequest(request.WalletFundedNotification, external_models.WalletFundedNotificationRequest{
			AccountID:     uint(businessID),
			Amount:        amount,
			Currency:      strings.Replace(currency, "ESCROW_", "", -1),
			TransactionID: transactionID,
		})
	}

	return walletBalance, nil
}

func DebitWallet(extReq request.ExternalRequest, db postgresql.Databases, amount float64, currency string, businessID int, creditEscrow string, transactionID string) (external_models.WalletBalance, error) {
	if currency == "" {
		currency = "NGN"
	}
	currency = strings.ToUpper(currency)

	if strings.ToLower(creditEscrow) == "yes" {
		currency = fmt.Sprintf("ESCROW_%v", currency)
		extReq.Logger.Info("debit-wallet", "creditEscrow is yes", "currency = ", currency, fmt.Sprintf("transaction with id %v", transactionID))
	}

	walletBalance, err := GetWalletBalanceByAccountIdAndCurrency(extReq, businessID, currency)
	if err != nil {
		walletBalance, err = CreateWalletBalance(extReq, businessID, currency, 0)
		if err != nil {
			return walletBalance, err
		}
		extReq.Logger.Info("debit-wallet-c", "new balance:", fmt.Sprintf("%v %v", currency, amount))
	}

	if amount > walletBalance.Available {
		return walletBalance, fmt.Errorf("insufficient wallet  balance")
	} else {
		availableBalance := walletBalance.Available - amount
		walletBalance, err = UpdateWalletBalance(extReq, walletBalance.ID, availableBalance)
		if err != nil {
			return walletBalance, err
		}

		extReq.Logger.Info("debit-wallet-u", "new balance:", fmt.Sprintf("%v %v", currency, availableBalance))
	}

	// walletDebitLog := models.WalletDebitLog{
	// 	AccountID: businessID,
	// 	Amount:    amount,
	// 	Currency:  currency,
	// }

	// walletDebitLog.CreateWalletDebitLog(db.Payment)
	extReq.SendExternalRequest(request.WalletDebitNotification, external_models.WalletDebitNotificationRequest{
		AccountID:     uint(businessID),
		Amount:        amount,
		Currency:      strings.Replace(currency, "ESCROW_", "", -1),
		TransactionID: transactionID,
	})

	return walletBalance, nil
}

func CreateWalletBalance(extReq request.ExternalRequest, accountID int, currency string, available float64) (external_models.WalletBalance, error) {
	walletItf, err := extReq.SendExternalRequest(request.CreateWalletBalance, external_models.CreateWalletRequest{
		AccountID: uint(accountID),
		Currency:  strings.ToUpper(currency),
		Available: available,
	})
	if err != nil {
		return external_models.WalletBalance{}, err
	}

	wallet, ok := walletItf.(external_models.WalletBalance)
	if !ok {
		return wallet, fmt.Errorf("response data format error")
	}

	return wallet, nil
}

func GetWalletBalanceByAccountIdAndCurrency(extReq request.ExternalRequest, accountID int, currency string) (external_models.WalletBalance, error) {
	walletItf, err := extReq.SendExternalRequest(request.GetWalletBalanceByAccountIDAndCurrency, external_models.GetWalletRequest{
		AccountID: uint(accountID),
		Currency:  strings.ToUpper(currency),
	})
	if err != nil {
		return external_models.WalletBalance{}, err
	}

	wallet, ok := walletItf.(external_models.WalletBalance)
	if !ok {
		return wallet, fmt.Errorf("response data format error")
	}

	return wallet, nil
}

func UpdateWalletBalance(extReq request.ExternalRequest, id uint, available float64) (external_models.WalletBalance, error) {
	walletItf, err := extReq.SendExternalRequest(request.UpdateWalletBalance, external_models.UpdateWalletRequest{
		ID:        id,
		Available: available,
	})
	if err != nil {
		return external_models.WalletBalance{}, err
	}

	wallet, ok := walletItf.(external_models.WalletBalance)
	if !ok {
		return wallet, fmt.Errorf("response data format error")
	}

	return wallet, nil
}

func UpdateTransactionAmountPaid(extReq request.ExternalRequest, transactionID string, amount float64, action string) (external_models.Transaction, error) {
	transactionItf, err := extReq.SendExternalRequest(request.UpdateTransactionAmountPaid, external_models.UpdateTransactionAmountPaidRequest{
		TransactionID: transactionID,
		Amount:        amount,
		Action:        action,
	})
	if err != nil {
		return external_models.Transaction{}, err
	}

	transaction, ok := transactionItf.(external_models.Transaction)
	if !ok {
		return transaction, fmt.Errorf("response data format error")
	}

	return transaction, nil
}

func CreateWalletTransaction(extReq request.ExternalRequest, senderAccountID, receiverAccountID int, senderAmount, receiverAmount float64, senderCurrency, receiverCurrency string, approved WalletTransactionApproved, firstApproval bool, secondApproval ...bool) (external_models.WalletTransaction, error) {
	data := external_models.CreateWalletTransactionRequest{
		SenderAccountID:   senderAccountID,
		ReceiverAccountID: receiverAccountID,
		SenderAmount:      senderAmount,
		ReceiverAmount:    receiverAmount,
		SenderCurrency:    senderCurrency,
		ReceiverCurrency:  receiverCurrency,
		Approved:          string(approved),
		FirstApproval:     firstApproval,
	}

	if len(secondApproval) > 0 {
		data.SecondApproval = &secondApproval[0]
	}

	walletTransactionInterface, err := extReq.SendExternalRequest(request.CreateWalletTransaction, data)
	if err != nil {
		extReq.Logger.Error(err.Error())
		return external_models.WalletTransaction{}, err
	}

	walletTransaction, ok := walletTransactionInterface.(external_models.WalletTransaction)
	if !ok {
		return walletTransaction, fmt.Errorf("response data format error")
	}

	return walletTransaction, nil
}

func CreateWalletHistory(extReq request.ExternalRequest, accountID int, reference string, amount, availableBalance float64, currency string, hType WalletHistoryType) (external_models.WalletHistory, error) {
	data := external_models.CreateWalletHistoryRequest{
		AccountID:        accountID,
		Reference:        reference,
		Amount:           amount,
		AvailableBalance: availableBalance,
		Currency:         strings.ToUpper(currency),
		Type:             string(hType),
	}

	walletHistoryInterface, err := extReq.SendExternalRequest(request.CreateWalletHistory, data)
	if err != nil {
		extReq.Logger.Error(err.Error())
		return external_models.WalletHistory{}, err
	}

	walletHistory, ok := walletHistoryInterface.(external_models.WalletHistory)
	if !ok {
		return walletHistory, fmt.Errorf("response data format error")
	}

	return walletHistory, nil
}

func SaveWalletHistory(extReq request.ExternalRequest, senderAccountID, receiverAccountID int, senderAmount, receiverAmount float64, senderCurrency, receiverCurrency string, senderAvailable, receiverAvailable float64) error {
	reference := utility.RandomString(10)

	_, err := CreateWalletHistory(extReq, senderAccountID, reference, senderAmount, senderAvailable, senderCurrency, WalletHistoryDebit)
	if err != nil {
		return err
	}

	_, err = CreateWalletHistory(extReq, receiverAccountID, reference, receiverAmount, receiverAvailable, receiverCurrency, WalletHistoryCredit)
	if err != nil {
		return err
	}
	return nil
}
