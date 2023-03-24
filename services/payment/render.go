package payment

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
)

func GetPaymentInvoiceService(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, paymentID string) (*template.Template, models.PaymentInvoiceData, int, error) {
	payment := models.Payment{PaymentID: paymentID}
	code, err := payment.GetPaymentByPaymentID(db.Payment)
	if err != nil {
		return &template.Template{}, models.PaymentInvoiceData{}, code, err
	}

	paymentInfo := models.PaymentInfo{PaymentID: paymentID}
	code, err = paymentInfo.GetPaymentInfoByPaymentID(db.Payment)
	if err != nil {
		return &template.Template{}, models.PaymentInvoiceData{}, code, err
	}

	transaction, _ := ListTransactionsByID(extReq, payment.TransactionID)
	buyer, _ := GetUserWithAccountID(extReq, transaction.Parties["buyer"].AccountID)
	seller, _ := GetUserWithAccountID(extReq, transaction.Parties["seller"].AccountID)
	brokerChargeBearer := transaction.Parties["broker_charge_bearer"]
	shippingChargeBearer := transaction.Parties["shipping_charge_bearer"]
	inspectionPeriod, _ := strconv.Atoi(transaction.InspectionPeriod)
	inspectionPeriodAsDate := ""
	if inspectionPeriod > 0 {
		t := time.Unix(int64(inspectionPeriod), 0)
		inspectionPeriodAsDate = t.Format("2006-01-02")
	}

	invoiceData := models.PaymentInvoiceData{
		Reference:        paymentInfo.Reference,
		PaymentID:        payment.PaymentID,
		TransactionID:    payment.TransactionID,
		TransactionType:  transaction.Type,
		Transaction:      transaction,
		Buyer:            buyer,
		Seller:           seller,
		InspectionPeriod: inspectionPeriodAsDate,
		ExpectedDelivery: transaction.DueDateFormatted,
		Title:            transaction.Title,
		Currency:         thisOrThatStr(payment.Currency, "NGN"),
		Amount:           payment.TotalAmount,
		EscrowCharge:     payment.EscrowCharge,
	}

	if strings.EqualFold(transaction.Type, "broker") && buyer.AccountID == uint(brokerChargeBearer.AccountID) {
		invoiceData.BrokerCharge = payment.BrokerCharge
	}

	if strings.EqualFold(transaction.Type, "broker") && buyer.AccountID == uint(shippingChargeBearer.AccountID) {
		invoiceData.ShippingFee = payment.ShippingFee
	}

	parsedTemplate, err := template.ParseFiles("./templates/invoice.html")
	if err != nil {
		return &template.Template{}, models.PaymentInvoiceData{}, http.StatusInternalServerError, err
	}
	invoiceData.TotalAmount = invoiceData.Amount + invoiceData.BrokerCharge + invoiceData.EscrowCharge + invoiceData.ShippingFee

	return parsedTemplate, invoiceData, http.StatusOK, nil
}
