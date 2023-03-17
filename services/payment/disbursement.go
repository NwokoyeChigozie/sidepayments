package payment

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/config"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/payment-ms/utility"
)

func WalletTransferService(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, req models.WalletTransferRequest) (string, int, error) {
	var (
		msg string
	)
	if req.RateID != 0 && req.InitialAmount == 0 {
		return msg, http.StatusBadRequest, fmt.Errorf("initial amount must be set if rate exist")
	}
	if req.RateID == 0 && req.InitialAmount != 0 {
		return msg, http.StatusBadRequest, fmt.Errorf("rate must be set if initial amount exist")
	}

	sender, err := GetUserWithAccountID(extReq, req.SenderAccountID)
	if err != nil {
		return msg, http.StatusInternalServerError, err
	}

	_, err = GetUserWithAccountID(extReq, req.RecipientAccountID)
	if err != nil {
		return msg, http.StatusInternalServerError, err
	}

	if !sender.CanMakeWithdrawal {
		return msg, http.StatusForbidden, fmt.Errorf("user withdrawal not enabled, please contact customer care")
	}

	if req.RateID != 0 && !sender.CanExchange {
		return msg, http.StatusForbidden, fmt.Errorf("user currency exchange not enabled, please contact customer care")
	}

	var (
		recipientCurrency      = strings.ToUpper(req.RecipientCurrency)
		senderCurrency         = strings.ToUpper(req.SenderCurrency)
		senderAvailableBalance float64
		amount                 float64
		initialCurrency        = strings.ReplaceAll(senderCurrency, "ESCROW_", "")
		finalCurrency          = strings.ReplaceAll(recipientCurrency, "ESCROW_", "")
		recipientAmount        float64
		escrowWallet           = "no"
	)

	senderWallet, err := GetWalletBalanceByAccountIdAndCurrency(extReq, req.SenderAccountID, senderCurrency)
	if err != nil {
		return msg, http.StatusBadRequest, fmt.Errorf("sender wallet does not exist")
	}
	senderAvailableBalance = senderWallet.Available

	if req.RateID != 0 && req.InitialAmount > 0 {
		amount = req.InitialAmount
	} else {
		amount = req.FinalAmount
	}

	if senderAvailableBalance < amount {
		return msg, http.StatusBadRequest, fmt.Errorf("insufficient balance")
	}

	if req.InitialAmount > 0 && req.RateID != 0 {
		if initialCurrency != finalCurrency {
			err := CreateExchangeTransaction(extReq, req.RecipientAccountID, req.RateID, req.InitialAmount, req.FinalAmount, ExchangeTransactionCompleted)
			if err != nil {
				return msg, http.StatusInternalServerError, err
			}
		}
	}

	if req.InitialAmount > 0 && req.RateID != 0 {
		rate, err := GetRateByID(extReq, req.RateID)
		if err != nil {
			return msg, http.StatusBadRequest, err
		}

		convertedAmount := (rate.Amount / req.InitialAmount) * amount
		recipientAmount = convertedAmount
	} else {
		recipientAmount = amount
	}

	senderWallet, debitStatus, err := DebitWallet(extReq, db, amount, req.SenderCurrency, req.SenderAccountID, escrowWallet, req.TransactionID)
	if err != nil {
		return msg, http.StatusInternalServerError, err
	} else if !debitStatus {
		return msg, http.StatusInternalServerError, fmt.Errorf("debit failed for sender")
	}

	if utility.InStringSlice(strings.ToUpper(req.SenderCurrency), []string{"NGN", "ESCROW_NGN", "USD", "GBP", "ESCROW_USD", "ESCROW_GBP"}) {
		if utility.InStringSlice(strings.ToUpper(req.SenderCurrency), []string{"NGN", "ESCROW_NGN"}) && amount >= config.GetConfig().ONLINE_PAYMENT.NairaThreshold {
			_, err := CreateWalletTransaction(extReq, req.SenderAccountID, req.RecipientAccountID, amount, recipientAmount, senderCurrency, recipientCurrency, WalletTransactionPending, false)
			if err != nil {
				return msg, http.StatusInternalServerError, err
			}
			return "Wallet transfer queued for approval by admin", http.StatusOK, nil
		}
	}

	receiverWallet, err := CreditWallet(extReq, db, recipientAmount, recipientCurrency, req.RecipientAccountID, req.Refund, escrowWallet, req.TransactionID)
	if err != nil {
		return msg, http.StatusInternalServerError, err
	}

	err = SaveWalletHistory(extReq, req.SenderAccountID, req.RecipientAccountID, amount, recipientAmount, senderCurrency, recipientCurrency, senderWallet.Available, receiverWallet.Available)
	if err != nil {
		extReq.Logger.Error("error saving wallet history: ", err.Error())
	}

	if req.TransactionID != "" {
		transaction, _ := ListTransactionsByID(extReq, req.TransactionID)
		transactionTitle := strings.Split(transaction.Title, ";")[0]
		action := "+"
		description := fmt.Sprintf("A sum of %v %v has been paid for this transaction", recipientCurrency, recipientAmount)
		if req.Refund {
			action = "-"
			description = fmt.Sprintf("A sum of %v %v has been deducted based on an excess payment being made on transaction %v", recipientCurrency, recipientAmount, transactionTitle)
		}

		extReq.SendExternalRequest(request.CreateActivityLog, external_models.CreateActivityLogRequest{
			TransactionID: req.TransactionID,
			Description:   description,
		})

		UpdateTransactionAmountPaid(extReq, req.TransactionID, amount, action)
	}

	return "Wallet transfer successful", http.StatusOK, nil
}
