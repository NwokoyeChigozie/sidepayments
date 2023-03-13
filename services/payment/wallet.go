package payment

import (
	"fmt"
	"strings"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/config"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
)

func CreditWallet(extReq request.ExternalRequest, db postgresql.Databases, amount float64, currency string, businessID int, isRefund bool, creditEscrow string, transactionID string) error {
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
			return err
		}
		extReq.Logger.Info("credit-wallet-c", "new balance:", fmt.Sprintf("%v %v", currency, amount))
	} else {
		availableBalance := walletBalance.Available + amount
		if isRefund {
			availableBalance += config.GetConfig().ONLINE_PAYMENT.DisbursementCharge
		}
		walletBalance, err = UpdateWalletBalance(extReq, walletBalance.ID, availableBalance)
		if err != nil {
			return err
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
			AccountID: uint(businessID),
			Amount:    amount,
		})
	}

	return nil
}

func DebitWallet(extReq request.ExternalRequest, db postgresql.Databases, amount float64, currency string, businessID int, creditEscrow string, transactionID string) (bool, error) {
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
			return false, err
		}
		extReq.Logger.Info("debit-wallet-c", "new balance:", fmt.Sprintf("%v %v", currency, amount))
	}

	if amount > walletBalance.Available {
		return false, nil
	} else {
		availableBalance := walletBalance.Available - amount
		walletBalance, err = UpdateWalletBalance(extReq, walletBalance.ID, availableBalance)
		if err != nil {
			return false, err
		}

		extReq.Logger.Info("debit-wallet-u", "new balance:", fmt.Sprintf("%v %v", currency, availableBalance))
	}

	walletDebitLog := models.WalletDebitLog{
		AccountID: businessID,
		Amount:    amount,
		Currency:  currency,
	}

	walletDebitLog.CreateWalletDebitLog(db.Payment)
	extReq.SendExternalRequest(request.WalletDebitNotification, external_models.WalletDebitNotificationRequest{
		AccountID:     uint(businessID),
		Amount:        amount,
		TransactionID: transactionID,
	})

	return true, nil
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
