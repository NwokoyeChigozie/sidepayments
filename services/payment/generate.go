package payment

import (
	"bytes"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/signintech/gopdf"
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

	data.Currency = currency
	data.TotalAmount = payment.EscrowCharge + brokerCharge + shippingCharge + payment.TotalAmount
	data.AmountPaid = payment.TotalAmount
	data.FeesPaid = payment.EscrowCharge
	if pdfLink != "" {
		data.PdfLink = pdfLink
	}

	return data
}

func GeneratePDFFromTemplate(templatePath string, data interface{}) (gopdf.GoPdf, error) {
	// Load the template
	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return gopdf.GoPdf{}, err
	}

	// Render the template to a string
	var renderedTemplate bytes.Buffer
	if err := tpl.Execute(&renderedTemplate, data); err != nil {
		return gopdf.GoPdf{}, err
	}

	// Create a new PDF document
	pdf := gopdf.GoPdf{}

	// Initialize the document
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	// Add a new page
	pdf.AddPage()

	// Set the font
	pdf.SetFont("Arial", "", 16)

	// Write the rendered template to the PDF
	pdf.Cell(nil, renderedTemplate.String())
	return pdf, nil
}

func GetPdfLink(extReq request.ExternalRequest, templatePath string, data interface{}) (string, error) {
	pdf, err := GeneratePDFFromTemplate(templatePath, data)
	if err != nil {
		return "", err
	}

	// Save the PDF to a bytes buffer
	var pdfBuf bytes.Buffer
	if err := pdf.Write(&pdfBuf); err != nil {
		return "", err
	}

	pdfItf, err := extReq.SendExternalRequest(request.UploadFile, external_models.UploadFileRequest{
		PlaceHolderName: "output.pdf",
		File:            pdfBuf.Bytes(),
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
