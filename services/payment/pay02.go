package payment

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/payment-ms/utility"
)

func ChargeCardInitService(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, req models.ChargeCardInitRequest) (string, int, error) {
	transaction, err := ListTransactionsByID(extReq, req.TransactionID)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	sellerParty, ok := transaction.Parties["seller"]
	if !ok {
		return "", http.StatusBadRequest, fmt.Errorf("transaction lacks a seller")
	}

	user, err := GetUserWithAccountID(extReq, sellerParty.AccountID)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("could not retrieve seller info: %v", err)
	}

	payment := models.Payment{PaymentID: req.PaymentID}
	code, err := payment.GetPaymentByPaymentID(db.Payment)
	if err != nil {
		return "", code, err
	}

	paymentCheck := models.Payment{TransactionID: payment.TransactionID}
	code, err = paymentCheck.GetPaymentByTransactionID(db.Payment)
	if err != nil {
		return "", code, fmt.Errorf("transaction has no payment record: %v", err.Error())
	}

	paymentCardInfo := models.PaymentCardInfo{AccountID: sellerParty.AccountID}
	code, err = paymentCardInfo.GetPaymentCardInfoByAccountID(db.Payment)
	if err != nil {
		return "", code, fmt.Errorf("recipient has no card details stored: %v", err.Error())
	}

	if paymentCardInfo.CardLifeTimeToken == "" {
		return "", http.StatusBadRequest, fmt.Errorf("user has no chargeable card")
	}

	rave := Rave{ExtReq: extReq}
	reference := fmt.Sprintf("VC%v", utility.RandomString(10))
	status, err := rave.ChargeCard(paymentCardInfo.CardLifeTimeToken, transaction.Currency, user.EmailAddress, reference, payment.TotalAmount)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	return status, http.StatusOK, nil
}
