package payment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/config"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/payment-ms/utility"
)

func GetStatusService(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, req models.GetStatusRequest) (string, int, error) {
	var (
		paymentInfo     = models.PaymentInfo{Reference: req.Reference}
		responseMessage = ""
		uri             = ""
		successPage     = ""
		msg             = ""
		escrowWallet    = "no"
		paymentGateway  = "rave"
		rave            = Rave{ExtReq: extReq}
		monnify         = Monnify{ExtReq: extReq}
		paymentChannelD = config.GetConfig().Slack.PaymentChannelID
	)

	code, err := paymentInfo.GetPaymentInfoByReference(db.Payment)
	if err != nil {
		return "error", code, fmt.Errorf("Payment data lacks a log record: %v", err.Error())
	}

	payment := models.Payment{PaymentID: paymentInfo.PaymentID}
	code, err = payment.GetPaymentByPaymentID(db.Payment)
	if err != nil {
		return msg, code, fmt.Errorf("Payment data lacks a payment record: %v", err.Error())
	}

	reqByte, err := json.Marshal(req)
	if err != nil {
		return responseMessage, http.StatusInternalServerError, err
	}

	callback := models.PaymentCallback{
		Log: string(reqByte),
	}
	err = callback.CreatePaymentCallback(db.Payment)
	if err != nil {
		return responseMessage, http.StatusInternalServerError, err
	}

	successPage = paymentInfo.RedirectUrl
	failPage := paymentInfo.FailUrl

	if paymentInfo.Status == "paid" {
		uri, err := HandleGetPaymentStatusPaid(c, extReq, payment, paymentInfo)
		if err != nil {
			return "error", http.StatusInternalServerError, err
		}

		if req.Headless {
			return "Transacton Already Paid", http.StatusOK, nil
		}

		c.Redirect(http.StatusMovedPermanently, uri)
		return "Transacton Already Paid", http.StatusOK, nil
	}

	var (
		gatewayStatus    = false
		businessID       int
		transactionID    string
		transactionTitle string
	)

	switch strings.ToLower(paymentGateway) {
	case "rave":
		gatewayStatus, _, err = rave.StatusV3(db, payment, paymentInfo, req.Reference)
		if err != nil {
			return "rave error", http.StatusInternalServerError, err
		}
	case "monnify":
		gatewayStatus, _, err = monnify.Status(req.Reference)
		if err != nil {
			return "monnify error", http.StatusInternalServerError, err
		}
	default:
		gatewayStatus, _, err = rave.StatusV3(db, payment, paymentInfo, req.Reference)
		if err != nil {
			return "rave error", http.StatusInternalServerError, err
		}
	}

	if gatewayStatus {
		paymentInfo.Status = "paid"
		err := paymentInfo.UpdateAllFields(db.Payment)
		if err != nil {
			return "error", http.StatusInternalServerError, err
		}

		payment.IsPaid = true
		payment.PaymentMadeAt = time.Now()
		if err != nil {
			return "payment update error", http.StatusInternalServerError, err
		}

		if payment.TransactionID != "" {
			transactionID = payment.TransactionID
			transaction, err := ListTransactionsByID(extReq, payment.TransactionID)
			if err != nil {
				return "error", http.StatusInternalServerError, err
			}
			transactionTitle = transaction.Title
			businessID = transaction.BusinessID
			if transaction.EscrowWallet != "" {
				escrowWallet = transaction.EscrowWallet
			}

			if transaction.MilestoneID != "" {
				extReq.SendExternalRequest(request.TransactionUpdateStatus, external_models.UpdateTransactionStatusRequest{
					AccountID:     businessID,
					TransactionID: transaction.TransactionID,
					MilestoneID:   transaction.MilestoneID,
					Status:        "ip",
				})

				for _, m := range transaction.Milestones {
					extReq.SendExternalRequest(request.TransactionUpdateStatus, external_models.UpdateTransactionStatusRequest{
						AccountID:     businessID,
						TransactionID: transaction.TransactionID,
						MilestoneID:   m.MilestoneID,
						Status:        "ip",
					})

					paymentUpdate := models.Payment{TransactionID: transaction.TransactionID}
					paymentUpdate.GetPaymentByTransactionID(db.Payment)
					paymentUpdate.IsPaid = true
					paymentUpdate.PaymentMadeAt = time.Now()
					paymentUpdate.UpdateAllFields(db.Payment)
				}

			} else {
				if transaction.Source == "transfer" {
					extReq.SendExternalRequest(request.BuyerSatisfied, external_models.OnlyTransactionIDRequiredRequest{
						TransactionID: transaction.TransactionID,
					})
				}
			}
			buyerParty := transaction.Parties["buyer"]
			businessEscrowCharge, err := getBusinessChargeWithBusinessIDAndCountry(extReq, transaction.BusinessID, transaction.Country.CountryCode)
			if err == nil {
				var (
					totalAmount       = payment.TotalAmount
					businessPerc, _   = strconv.ParseFloat(businessEscrowCharge.BusinessCharge, 64)
					vesicashCharge, _ = strconv.ParseFloat(businessEscrowCharge.VesicashCharge, 64)
					amountOne         = utility.PercentageOf(totalAmount, businessPerc)
					amountTwo         = utility.PercentageOf(totalAmount, vesicashCharge)
				)

				//credit vesicash
				_, err = CreditWallet(extReq, db, amountTwo, transaction.Currency, 1, false, "no", transaction.TransactionID)
				if err != nil {
					return "error", http.StatusInternalServerError, err
				}
				amountOne = totalAmount - amountOne
				_, err = CreditWallet(extReq, db, amountOne, transaction.Currency, buyerParty.AccountID, false, escrowWallet, transaction.TransactionID)
				if err != nil {
					return "error", http.StatusInternalServerError, err
				}

			} else {
				_, err = CreditWallet(extReq, db, payment.TotalAmount, transaction.Currency, buyerParty.AccountID, false, escrowWallet, transaction.TransactionID)
				if err != nil {
					return "error", http.StatusInternalServerError, err
				}
			}

			err = SlackNotify(paymentChannelD, `
 					Payment Status For Transaction #`+payment.TransactionID+`
                     Environment: `+config.GetConfig().App.Name+`
                     Payment ID: `+paymentInfo.PaymentID+`
                     Amount: `+fmt.Sprintf("%v %v", payment.Currency, payment.TotalAmount)+`
                     Status: SUCCESSFUL
 				`)
			if err != nil && !extReq.Test {
				extReq.Logger.Error("error sending notification to slack: ", err.Error())
			}

		} else {
			if req.FundWallet {
				_, err = CreditWallet(extReq, db, payment.TotalAmount, payment.Currency, int(payment.AccountID), false, escrowWallet, "")
				if err != nil {
					return "error", http.StatusInternalServerError, err
				}
				beneficiary, _ := GetUserWithAccountID(extReq, int(payment.AccountID))
				err = SlackNotify(paymentChannelD, `
 					Wallet Disbursement/Funding For Customer #`+fmt.Sprintf("%v", payment.AccountID)+`
                     Environment: `+config.GetConfig().App.Name+`
                     Account ID: `+fmt.Sprintf("%v", payment.AccountID)+`
                     Beneficiary Name: `+fmt.Sprintf("%v %v", beneficiary.Firstname, beneficiary.Lastname)+`
                     Amount: `+fmt.Sprintf("%v %v", payment.Currency, payment.TotalAmount)+`
                     Status: Success
 				`)
				if err != nil && !extReq.Test {
					extReq.Logger.Error("error sending notification to slack: ", err.Error())
				}

				err = SlackNotify(paymentChannelD, `
 					Payment Status For Headless Payment #`+payment.PaymentID+`
 					Environment: `+config.GetConfig().App.Name+`
 					Payment ID: `+payment.PaymentID+`
 					Amount: `+fmt.Sprintf("%v %v", payment.Currency, payment.TotalAmount)+`
 					Status: SUCCESSFUL
 				`)
				if err != nil && !extReq.Test {
					extReq.Logger.Error("error sending notification to slack: ", err.Error())
				}
			}
		}

		pdfLink, _ := utility.URLDecode(utility.GenerateGroupByURL(c, fmt.Sprintf("payment/invoice/%v", payment.PaymentID), map[string]string{}))
		businessProfileData, _ := GetBusinessProfileByAccountID(extReq, extReq.Logger, businessID)
		if successPage != "" {
			utility.AddQueryParam(&successPage, "invoice", pdfLink)
			utility.AddQueryParam(&successPage, "transaction_id", transactionID)
			utility.AddQueryParam(&successPage, "reference", paymentInfo.Reference)
		} else {
			if businessProfileData.RedirectUrl != "" {
				uri = businessProfileData.RedirectUrl
				utility.AddQueryParam(&uri, "source", "redirect_url")
			} else {
				uri = config.GetConfig().App.SiteUrl + "/paylink/success"
				utility.AddQueryParam(&uri, "source", "default")
			}
			utility.AddQueryParam(&uri, "invoice", pdfLink)
			utility.AddQueryParam(&uri, "transaction_id", transactionID)
			utility.AddQueryParam(&uri, "reference", paymentInfo.Reference)

			if businessProfileData.Webhook_uri != "" {
				InitWebhook(extReq, db, businessProfileData.Webhook_uri, "payment.success", map[string]interface{}{
					"transaction_title": transactionTitle,
					"transaction_id":    transactionID,
					"payment_status":    "success",
				}, businessProfileData.AccountID)
			}

		}

		if req.Headless {
			return "Transaction payment successfully confirmed", http.StatusOK, nil
		}
		c.Redirect(http.StatusMovedPermanently, uri)
		return "Transaction payment successfully confirmed", http.StatusOK, nil
	} else {
		paymentInfo.Status = "failed"
		err := paymentInfo.UpdateAllFields(db.Payment)
		if err != nil {
			return "error", http.StatusInternalServerError, err
		}

		if payment.TransactionID != "" {
			transactionID = payment.TransactionID
			transaction, err := ListTransactionsByID(extReq, payment.TransactionID)
			if err != nil {
				return "error", http.StatusInternalServerError, err
			}
			transactionTitle = transaction.Title
			businessID = transaction.BusinessID
			businessProfileData, _ := GetBusinessProfileByAccountID(extReq, extReq.Logger, businessID)
			if businessProfileData.Webhook_uri != "" {
				InitWebhook(extReq, db, businessProfileData.Webhook_uri, "payment.failed", map[string]interface{}{
					"transaction_title": transactionTitle,
					"transaction_id":    transactionID,
					"payment_status":    "failed",
				}, businessProfileData.AccountID)
			}

			err = SlackNotify(paymentChannelD, `
 					Payment Status For Transaction #`+payment.TransactionID+`
                     Environment: `+config.GetConfig().App.Name+`
                     Payment ID: `+paymentInfo.PaymentID+`
                     Amount: `+fmt.Sprintf("%v %v", payment.Currency, payment.TotalAmount)+`
                     Status: FAILED
                     
 				`)
			if err != nil && !extReq.Test {
				extReq.Logger.Error("error sending notification to slack: ", err.Error())
			}
		} else {
			err = SlackNotify(paymentChannelD, `
 					Payment Status For Headless Payment #`+payment.PaymentID+`
                     Environment: `+config.GetConfig().App.Name+`
                     Payment ID: `+paymentInfo.PaymentID+`
                     Amount: `+fmt.Sprintf("%v %v", payment.Currency, payment.TotalAmount)+`
                     Status: FAILED
                     
 				`)
			if err != nil && !extReq.Test {
				extReq.Logger.Error("error sending notification to slack: ", err.Error())
			}
		}
	}

	if !req.Headless {
		uri = failPage
		if uri == "" {
			uri = config.GetConfig().App.Url + "/payment/pay/failed"
		}
		c.Redirect(http.StatusMovedPermanently, uri)
		return "Transaction payment failed", http.StatusBadRequest, nil
	}

	return "Transaction payment failed", http.StatusBadRequest, nil
}

func GetPaymentStatusService(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, req models.GetPaymentStatusRequest) (string, string, int, error) {
	var (
		paymentInfo     = models.PaymentInfo{Reference: req.Reference}
		uri             = ""
		msg             = ""
		paymentGateway  = "rave"
		rave            = Rave{ExtReq: extReq}
		monnify         = Monnify{ExtReq: extReq}
		paymentChannelD = config.GetConfig().Slack.PaymentChannelID
	)

	code, err := paymentInfo.GetPaymentInfoByReference(db.Payment)
	if err != nil {
		return uri, "error", code, fmt.Errorf("payment data lacks a log record: %v", err.Error())
	}

	payment := models.Payment{PaymentID: paymentInfo.PaymentID}
	code, err = payment.GetPaymentByPaymentID(db.Payment)
	if err != nil {
		return uri, msg, code, fmt.Errorf("payment data lacks a payment record: %v", err.Error())
	}

	reqByte, err := json.Marshal(req)
	if err != nil {
		return uri, msg, http.StatusInternalServerError, err
	}

	callback := models.PaymentCallback{
		Log: string(reqByte),
	}
	err = callback.CreatePaymentCallback(db.Payment)
	if err != nil {
		return uri, msg, http.StatusInternalServerError, err
	}

	if paymentInfo.Status == "paid" {
		uri, err := HandleGetPaymentStatusPaid(c, extReq, payment, paymentInfo)
		if err != nil {
			return uri, "error", http.StatusInternalServerError, err
		}
		return uri, "Transacton Already Paid", http.StatusOK, nil
	}

	var (
		gatewayStatus    = false
		successPage      = paymentInfo.RedirectUrl
		businessID       int
		escrowCharge     float64
		escrowWallet     = "no"
		transactionID    string
		transactionTitle string
	)

	switch strings.ToLower(paymentGateway) {
	case "rave":
		gatewayStatus, _, err = rave.StatusV3(db, payment, paymentInfo, req.Reference)
		if err != nil {
			return uri, "rave error", http.StatusInternalServerError, err
		}
	case "monnify":
		gatewayStatus, _, err = monnify.Status(req.Reference)
		if err != nil {
			return uri, "monnify error", http.StatusInternalServerError, err
		}
	default:
		gatewayStatus, _, err = rave.StatusV3(db, payment, paymentInfo, req.Reference)
		if err != nil {
			return uri, "rave error", http.StatusInternalServerError, err
		}
	}

	if gatewayStatus {
		paymentInfo.Status = "paid"
		err := paymentInfo.UpdateAllFields(db.Payment)
		if err != nil {
			return uri, "error", http.StatusInternalServerError, err
		}

		if payment.TransactionID != "" {
			transactionID = payment.TransactionID
			transaction, err := ListTransactionsByID(extReq, payment.TransactionID)
			if err != nil {
				return uri, "error", http.StatusInternalServerError, err
			}
			transactionTitle = transaction.Title
			businessID = transaction.BusinessID
			escrowCharge = transaction.EscrowCharge
			if transaction.EscrowWallet != "" {
				escrowWallet = transaction.EscrowWallet
			}

			// credit vesicash
			_, err = CreditWallet(extReq, db, escrowCharge, transaction.Currency, 1, false, "no", transaction.TransactionID)
			if err != nil {
				return uri, "error", http.StatusInternalServerError, err
			}
			buyerAmount := payment.TotalAmount - escrowCharge
			_, err = CreditWallet(extReq, db, buyerAmount, transaction.Currency, businessID, false, escrowWallet, transaction.TransactionID)
			if err != nil {
				return uri, "error", http.StatusInternalServerError, err
			}

			extReq.SendExternalRequest(request.CreateActivityLog, external_models.CreateActivityLogRequest{
				TransactionID: transaction.TransactionID,
				Description:   fmt.Sprintf("A sum of %v has been paid for this transaction", buyerAmount),
			})

			err = SlackNotify(paymentChannelD, `
					Payment Status For Transaction #`+payment.TransactionID+`
                    Environment: `+config.GetConfig().App.Name+`
                    Payment ID: `+paymentInfo.PaymentID+`
                    Amount: `+payment.Currency+` `+fmt.Sprintf("%v", payment.TotalAmount)+`
                    Status: SUCCESSFUL
			`)
			if err != nil && !extReq.Test {
				extReq.Logger.Error("error sending notification to slack: ", err.Error())
			}
		} else {
			if req.FundWallet {
				_, err = CreditWallet(extReq, db, payment.TotalAmount, payment.Currency, int(payment.AccountID), false, escrowWallet, "")
				if err != nil {
					return uri, "error", http.StatusInternalServerError, err
				}
				businessID = int(payment.BusinessID)
				user, _ := GetUserWithAccountID(extReq, int(payment.AccountID))
				err = SlackNotify(paymentChannelD, `
					Wallet Funding For Customer #`+fmt.Sprintf("%v", payment.AccountID)+` 
                    Environment: `+config.GetConfig().App.Name+`
                    Account ID: `+fmt.Sprintf("%v", payment.AccountID)+`
                    Beneficiary Name: `+fmt.Sprintf("%v %v", user.Firstname, user.Lastname)+`
                    Amount: `+payment.Currency+` `+fmt.Sprintf("%v", payment.TotalAmount)+`
                    Status: Success
				`)
				if err != nil && !extReq.Test {
					extReq.Logger.Error("error sending notification to slack: ", err.Error())
				}

				extReq.SendExternalRequest(request.PaymentInvoiceNotification, external_models.PaymentInvoiceNotificationRequest{
					Reference:                 req.Reference,
					PaymentID:                 payment.PaymentID,
					TransactionType:           "",
					TransactionID:             payment.TransactionID,
					Buyer:                     int(payment.AccountID),
					Seller:                    int(payment.BusinessID),
					InspectionPeriodFormatted: "",
					ExpectedDelivery:          "",
					Title:                     "",
					Currency:                  payment.Currency,
					Amount:                    payment.TotalAmount,
					EscrowCharge:              payment.EscrowCharge,
					BrokerCharge:              payment.BrokerCharge,
				})

				err = SlackNotify(paymentChannelD, `
					Payment Status For Headless Payment #`+payment.PaymentID+`                          
					Environment: `+config.GetConfig().App.Name+`
					Payment ID: `+payment.PaymentID+`
					Amount: `+payment.Currency+` `+fmt.Sprintf("%v", payment.TotalAmount)+`
					Status: SUCCESSFUL
				`)
				if err != nil && !extReq.Test {
					extReq.Logger.Error("error sending notification to slack: ", err.Error())
				}

			}

		}

		payment.IsPaid = true
		payment.PaymentMadeAt = time.Now()
		if err != nil {
			return uri, "payment update error", http.StatusInternalServerError, err
		}

		pdfLink, _ := utility.URLDecode(utility.GenerateGroupByURL(c, fmt.Sprintf("payment/invoice/%v", payment.PaymentID), map[string]string{}))
		businessProfileData, _ := GetBusinessProfileByAccountID(extReq, extReq.Logger, businessID)
		if successPage != "" {
			utility.AddQueryParam(&successPage, "invoice", pdfLink)
			utility.AddQueryParam(&successPage, "transaction_id", transactionID)
			utility.AddQueryParam(&successPage, "reference", paymentInfo.Reference)
		} else {
			if businessProfileData.RedirectUrl != "" {
				uri = businessProfileData.RedirectUrl
				utility.AddQueryParam(&uri, "source", "redirect_url")
			} else {
				uri = config.GetConfig().App.SiteUrl + "/paylink/success"
				utility.AddQueryParam(&uri, "source", "default")
			}
			utility.AddQueryParam(&uri, "invoice", pdfLink)
			utility.AddQueryParam(&uri, "transaction_id", transactionID)
			utility.AddQueryParam(&uri, "reference", paymentInfo.Reference)

			if businessProfileData.Webhook_uri != "" {
				InitWebhook(extReq, db, businessProfileData.Webhook_uri, "payment.success", map[string]interface{}{
					"transaction_title": transactionTitle,
					"transaction_id":    transactionID,
					"payment_status":    "success",
				}, businessProfileData.AccountID)
			}

		}
	}

	return uri, "Transaction payment successfully confirmed", http.StatusOK, nil
}

func HandleGetPaymentStatusPaid(c *gin.Context, extReq request.ExternalRequest, payment models.Payment, paymentInfo models.PaymentInfo) (string, error) {
	var (
		businessID    = 0
		transactionID = ""
		successPage   = paymentInfo.RedirectUrl
		uri           = ""
	)
	if payment.TransactionID != "" {
		transactionID = payment.TransactionID
		transaction, err := ListTransactionsByID(extReq, payment.TransactionID)
		if err != nil {
			return uri, err
		}
		businessID = transaction.BusinessID
	} else {
		businessID = int(payment.BusinessID)
	}

	pdfLink, err := utility.URLDecode(utility.GenerateGroupByURL(c, fmt.Sprintf("payment/invoice/%v", payment.PaymentID), map[string]string{}))
	if err != nil {
		return uri, fmt.Errorf("error generating invoice link")
	}

	businessProfileData, _ := GetBusinessProfileByAccountID(extReq, extReq.Logger, businessID)

	if successPage != "" {
		utility.AddQueryParam(&successPage, "invoice", pdfLink)
		utility.AddQueryParam(&successPage, "transaction_id", transactionID)
		utility.AddQueryParam(&successPage, "reference", paymentInfo.Reference)
	} else {
		if businessProfileData.RedirectUrl != "" {
			uri = businessProfileData.RedirectUrl
			utility.AddQueryParam(&uri, "source", "redirect_url")
		} else {
			uri = config.GetConfig().App.SiteUrl + "/paylink/success"
			utility.AddQueryParam(&uri, "source", "default")
		}
		utility.AddQueryParam(&uri, "invoice", pdfLink)
		utility.AddQueryParam(&uri, "transaction_id", transactionID)
		utility.AddQueryParam(&uri, "reference", paymentInfo.Reference)
	}
	return uri, nil
}
