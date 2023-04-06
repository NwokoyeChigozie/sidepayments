package payment

import (
	"bytes"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"

	wkhtmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/models"
)

type PdfData struct {
	TransactionType        string
	Currency               string
	AmountPaid             float64
	TotalAmount            float64
	Title                  string
	ExpectedDelivery       string
	InspectionPeriodAsDate string
	BuyerEmailAddress      string
	SellerEmailAddress     string
	ShippingFee            float64
	BrokerCharge           float64
	FeesPaid               float64
	TransactionID          string
	PdfLink                string
}

func NewPdfData(extReq request.ExternalRequest, transaction external_models.TransactionByID, payment models.Payment, reference, pdfLink string) PdfData {
	var (
		data = PdfData{
			Title:                  "-",
			ExpectedDelivery:       "-",
			InspectionPeriodAsDate: "-",
			BuyerEmailAddress:      "-",
			SellerEmailAddress:     "-",
			PdfLink:                "#",
		}
		brokerCharge   float64 = 0
		shippingCharge float64 = 0
		currency               = "NGN"
	)

	if transaction.ID != 0 {
		data.TransactionID = transaction.TransactionID
		buyer := transaction.Parties["buyer"]
		seller := transaction.Parties["seller"]
		brokerChargeBearer := transaction.Parties["broker_charge_bearer"]
		shippingChargeBearer := transaction.Parties["shipping_charge_bearer"]
		if transaction.Type == "broker" && buyer.AccountID == brokerChargeBearer.AccountID {
			brokerCharge = payment.BrokerCharge
		}
		if transaction.Type == "broker" && buyer.AccountID == shippingChargeBearer.AccountID {
			shippingCharge = payment.ShippingFee
		}
		if transaction.Currency != "" {
			currency = strings.ToUpper(transaction.Currency)
		}
		if transaction.Title != "" {
			data.Title = transaction.Title
		}

		if transaction.Source != "transfer" && transaction.DueDate != "" {
			data.ExpectedDelivery = transaction.DueDate
		}

		inspectionPeriod, _ := strconv.Atoi(transaction.InspectionPeriod)
		if inspectionPeriod > 0 {
			t := time.Unix(int64(inspectionPeriod), 0)
			data.InspectionPeriodAsDate = t.Format("2006-01-02")
		}

		buyerUser, _ := GetUserWithAccountID(extReq, buyer.AccountID)
		sellerUser, _ := GetUserWithAccountID(extReq, seller.AccountID)

		if buyerUser.EmailAddress != "" {
			data.BuyerEmailAddress = buyerUser.EmailAddress
		}
		if sellerUser.EmailAddress != "" {
			data.SellerEmailAddress = sellerUser.EmailAddress
		}

		if transaction.Type == "broker" {
			data.BrokerCharge = payment.BrokerCharge
		}
	}

	data.Currency = thisOrThatStr(currency, payment.Currency)
	data.TotalAmount = payment.EscrowCharge + brokerCharge + shippingCharge + payment.TotalAmount
	data.AmountPaid = payment.TotalAmount
	data.FeesPaid = payment.EscrowCharge
	if pdfLink != "" {
		data.PdfLink = pdfLink
	}

	return data
}

func GeneratePDFFromTemplate(templatePath string, data interface{}) ([]byte, error) {

	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	var renderedTemplate bytes.Buffer
	if err := tpl.Execute(&renderedTemplate, data); err != nil {
		return nil, err
	}

	html := renderedTemplate.String()

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(html)))

	if err := pdfg.Create(); err != nil {
		return nil, err
	}

	return pdfg.Bytes(), nil
}

func GetPdfLink(extReq request.ExternalRequest, templatePath string, data interface{}) (string, error) {
	pdf, err := GeneratePDFFromTemplate(templatePath, data)
	if err != nil {
		return "", err
	}

	pdfItf, err := extReq.SendExternalRequest(request.UploadFile, external_models.UploadFileRequest{
		PlaceHolderName: "output.pdf",
		File:            pdf,
	})
	if err != nil {
		return "", err
	}

	fileData, ok := pdfItf.(external_models.UploadFileResponseData)
	if !ok {
		return "", fmt.Errorf("response data format error")
	}

	return fileData.FileUrl, nil
}
