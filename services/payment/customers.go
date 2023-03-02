package payment

import (
	"net/http"

	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
)

func GetCustomerPaymentsService(extReq request.ExternalRequest, db postgresql.Databases, businessID int) (float64, int, error) {
	var (
		total       float64
		allPayments []models.Payment
	)
	users, err := GetUsersByBusinessID(extReq, businessID)
	if err != nil {
		return total, http.StatusInternalServerError, err
	}

	for _, user := range users {
		payment := models.Payment{AccountID: int64(user.AccountID), IsPaid: true}
		payments, err := payment.GetAllPaymentsByAccountIDAndIsPaidAndPaymentMadeAtNotNull(db.Payment)
		if err != nil {
			return total, http.StatusInternalServerError, err
		} else {
			allPayments = append(allPayments, payments...)
		}
	}

	for _, payment := range allPayments {
		total += payment.TotalAmount
	}

	return total, http.StatusOK, nil
}

func GetCustomerCardService(extReq request.ExternalRequest, db postgresql.Databases, accountID uint) (models.CardResponse, int, error) {
	var (
		paymentCardInfo = models.PaymentCardInfo{AccountID: int(accountID)}
	)

	code, err := paymentCardInfo.GetPaymentCardInfoByAccountID(db.Payment)
	if err != nil {
		return models.CardResponse{}, code, err
	}

	return models.CardResponse{
		AccountID:   int(accountID),
		Card:        paymentCardInfo.LastFourDigits,
		ExpiryMonth: paymentCardInfo.CcExpiryMonth,
		ExpiryYear:  paymentCardInfo.CcExpiryYear,
	}, http.StatusOK, nil
}

func GetCustomerCardsByBusinessIDService(extReq request.ExternalRequest, db postgresql.Databases, businessID int) ([]models.CardResponse, int, error) {
	var (
		cards           = []models.CardResponse{}
		paymentCardInfo = models.PaymentCardInfo{}
		accountIds      = []int{}
	)

	users, err := GetUsersByBusinessID(extReq, businessID)
	if err != nil {
		return cards, http.StatusInternalServerError, err
	}

	for _, user := range users {
		accountIds = append(accountIds, int(user.AccountID))
	}

	paymentCardInfos, err := paymentCardInfo.GetAllPaymentCardInfosByAccountIDs(db.Payment, accountIds)
	if err != nil {
		return cards, http.StatusInternalServerError, err
	}

	for _, p := range paymentCardInfos {
		cards = append(cards, models.CardResponse{
			AccountID:   p.AccountID,
			Card:        p.LastFourDigits,
			ExpiryMonth: p.CcExpiryMonth,
			ExpiryYear:  p.CcExpiryYear,
		})
	}

	return cards, http.StatusOK, nil
}
