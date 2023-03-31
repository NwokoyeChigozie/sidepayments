package test_payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/external/mocks/auth_mocks"
	"github.com/vesicash/payment-ms/external/mocks/transactions_mocks"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/config"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/controller/payment"
	"github.com/vesicash/payment-ms/pkg/middleware"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	tst "github.com/vesicash/payment-ms/tests"
	"github.com/vesicash/payment-ms/utility"
)

func TestGetStatus(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := postgresql.Connection()
	var (
		muuid, _  = uuid.NewV4()
		accountID = uint(utility.GetRandomNumbersInRange(1000000000, 9999999999))
		testUser  = external_models.User{
			ID:           uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
			AccountID:    accountID,
			EmailAddress: fmt.Sprintf("testuser%v@qa.team", muuid.String()),
			PhoneNumber:  fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
			AccountType:  "individual",
			Firstname:    "test",
			Lastname:     "user",
			Username:     fmt.Sprintf("test_username%v", muuid.String()),
		}
	)

	auth_mocks.User = &testUser
	auth_mocks.BusinessProfile = &external_models.BusinessProfile{
		ID:        uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		AccountID: int(testUser.AccountID),
		Country:   "NG",
		Currency:  "NGN",
	}
	auth_mocks.ValidateAuthorizationRes = &external_models.ValidateAuthorizationDataModel{
		Status:  true,
		Message: "authorized",
		Data:    testUser,
	}
	auth_mocks.UserProfile = &external_models.UserProfile{
		ID:        uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		AccountID: int(testUser.AccountID),
		Country:   "NG",
		Currency:  "NGN",
	}

	auth_mocks.BusinessCharge = &external_models.BusinessCharge{
		ID:                  uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		BusinessId:          int(testUser.AccountID),
		Country:             "NG",
		Currency:            "NGN",
		BusinessCharge:      "0",
		VesicashCharge:      "2.5",
		ProcessingFee:       "0",
		PaymentGateway:      "rave",
		DisbursementGateway: "rave_momo",
		ProcessingFeeMode:   "fixed",
	}

	auth_mocks.Country = &external_models.Country{
		ID:           uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		Name:         "nigeria",
		CountryCode:  "NG",
		CurrencyCode: "NGN",
	}

	var (
		transactionID = utility.RandomString(20)
		partiesID     = utility.RandomString(20)
		milestoneID   = utility.RandomString(20)
		reference     = utility.RandomString(20)
		reference2    = utility.RandomString(20)
	)

	transactions_mocks.ListTransactionsByIDObj = &external_models.TransactionByID{
		ID:               uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		TransactionID:    transactionID,
		PartiesID:        partiesID,
		MilestoneID:      milestoneID,
		Title:            "title",
		Type:             "milestone",
		Description:      "description",
		Amount:           1000,
		Status:           "draft",
		Quantity:         1,
		InspectionPeriod: "2023-03-10",
		DueDate:          "1678924800",
		ShippingFee:      0,
		GracePeriod:      "1679097600",
		Currency:         "NGN",
		BusinessID:       int(testUser.AccountID),
		IsPaylinked:      false,
		Source:           "api",
		TransUssdCode:    74950,
		Recipients: []external_models.MileStoneRecipient{
			{
				AccountID:    9489042479,
				Amount:       500,
				EmailAddress: "test.qa.team",
				PhoneNumber:  "+23456789776789",
			},
		},
		EscrowWallet: "yes",
		Parties: map[string]external_models.TransactionParty{
			"buyer": {
				ID:                   uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
				TransactionPartiesID: partiesID,
				TransactionID:        transactionID,
				AccountID:            int(testUser.AccountID),
				Role:                 "buyer",
				Status:               "draft",
				RoleCapabilities: &map[string]interface{}{"approve": true,
					"can_receive":  false,
					"can_view":     true,
					"mark_as_done": false,
				},
			},
			"seller": {
				ID:                   uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
				TransactionPartiesID: partiesID,
				TransactionID:        transactionID,
				AccountID:            int(testUser.AccountID),
				Role:                 "seller",
				Status:               "draft",
				RoleCapabilities: &map[string]interface{}{"approve": true,
					"can_receive":  false,
					"can_view":     true,
					"mark_as_done": false,
				},
			},
		},
		Members: []external_models.PartyResponse{
			{
				PartyID:     utility.GetRandomNumbersInRange(1000000000, 9999999999),
				AccountID:   int(testUser.AccountID),
				AccountName: "name",
				Email:       "email@email.com",
				PhoneNumber: "+2340905964039",
				Role:        "buyer",
				Status:      "draft",
				AccessLevel: external_models.PartyAccessLevel{
					CanView:    true,
					CanReceive: false,
					MarkAsDone: false,
					Approve:    true,
				},
			},
			{
				PartyID:     utility.GetRandomNumbersInRange(1000000000, 9999999999),
				AccountID:   int(testUser.AccountID),
				AccountName: "name",
				Email:       "email@email.com",
				PhoneNumber: "+2340905964039",
				Role:        "seller",
				Status:      "draft",
				AccessLevel: external_models.PartyAccessLevel{
					CanView:    true,
					CanReceive: false,
					MarkAsDone: false,
					Approve:    true,
				},
			},
		},
		TotalAmount: 2000,
		Milestones: []external_models.MilestonesResponse{
			{
				Index:            1,
				MilestoneID:      utility.RandomString(20),
				Title:            "milestone title",
				Amount:           1000,
				Status:           "draft",
				InspectionPeriod: "2023-03-10",
				DueDate:          "1678924800",
				Recipients: []external_models.MilestonesRecipientResponse{
					{
						AccountID:   9489042479,
						AccountName: "name",
						Amount:      500,
						Email:       "test.qa.team",
						PhoneNumber: "+23456789776789",
					},
				},
			},
			{
				Index:            2,
				MilestoneID:      utility.RandomString(20),
				Title:            "milestone title 1",
				Amount:           1000,
				Status:           "draft",
				InspectionPeriod: "2023-03-10",
				DueDate:          "1678924800",
				Recipients: []external_models.MilestonesRecipientResponse{
					{
						AccountID:   9489042479,
						AccountName: "name",
						Amount:      500,
						Email:       "test.qa.team",
						PhoneNumber: "+23456789776789",
					},
				},
			},
		},
		IsDisputed: false,
	}

	paymentData := models.Payment{
		ID:               int64(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		TransactionID:    transactionID,
		PaymentID:        utility.RandomString(20),
		TotalAmount:      200,
		EscrowCharge:     0,
		IsPaid:           false,
		AccountID:        int64(accountID),
		BusinessID:       int64(accountID),
		Currency:         "NGN",
		DisburseCurrency: "NGN",
	}

	err := paymentData.CreatePayment(db.Payment)
	if err != nil {
		t.Fatal("errpr creating payment: " + err.Error())
	}
	paymentData2 := models.Payment{
		ID:               int64(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		TransactionID:    transactionID,
		PaymentID:        utility.RandomString(20),
		TotalAmount:      200,
		EscrowCharge:     0,
		IsPaid:           false,
		AccountID:        int64(accountID),
		BusinessID:       int64(accountID),
		Currency:         "NGN",
		DisburseCurrency: "NGN",
	}

	err = paymentData2.CreatePayment(db.Payment)
	if err != nil {
		t.Fatal("errpr creating payment: " + err.Error())
	}

	paymentCardInfo := models.PaymentCardInfo{
		AccountID:         int(testUser.AccountID),
		PaymentID:         paymentData.PaymentID,
		CcExpiryMonth:     "12",
		CcExpiryYear:      "2025",
		LastFourDigits:    "4566",
		Brand:             "VISA",
		IssuingCountry:    "nigeria",
		CardToken:         "",
		CardLifeTimeToken: "902be0d9be02be09b",
		Payload:           "",
	}
	err = paymentCardInfo.CreatePaymentCardInfo(db.Payment)
	if err != nil {
		t.Fatal("error creating payment card info: " + err.Error())
	}

	paymentAccount1 := models.PaymentAccount{
		PaymentAccountID: reference,
		TransactionID:    transactionID,
		PaymentID:        paymentData.PaymentID,
		IsUsed:           true,
		ExpiresAfter:     strconv.Itoa(int(time.Now().Add(72 * time.Hour).Unix())),
		BusinessID:       strconv.Itoa(int(testUser.AccountID)),
		AccountNumber:    config.GetConfig().Rave.MerchantId,
		BankCode:         "flutterwave",
		BankName:         "flutterwave",
		Status:           "ACTIVE",
	}
	err = paymentAccount1.CreatePaymentAccount(db.Payment)
	if err != nil {
		t.Fatal("error creating payment account: " + err.Error())
	}

	paymentAccount2 := models.PaymentAccount{
		PaymentID:        paymentData2.PaymentID,
		PaymentAccountID: reference2,
		IsUsed:           true,
		ExpiresAfter:     strconv.Itoa(int(time.Now().Add(72 * time.Hour).Unix())),
		BusinessID:       strconv.Itoa(int(testUser.AccountID)),
		AccountNumber:    "7727632865",
		BankCode:         "221",
		BankName:         "vesicash bank",
		Status:           "ACTIVE",
	}
	err = paymentAccount2.CreatePaymentAccount(db.Payment)
	if err != nil {
		t.Fatal("error creating payment account: " + err.Error())
	}

	paymentInfo1 := models.PaymentInfo{
		PaymentID:   paymentData.PaymentID,
		Reference:   reference,
		Status:      "pending",
		Gateway:     "rave",
		RedirectUrl: "",
	}

	err = paymentInfo1.CreatePaymentInfo(db.Payment)
	if err != nil {
		t.Fatal("error creating payment info: " + err.Error())
	}

	paymentInfo2 := models.PaymentInfo{
		PaymentID:   paymentData2.PaymentID,
		Reference:   reference2,
		Status:      "pending",
		Gateway:     "monnify",
		RedirectUrl: "",
	}

	err = paymentInfo2.CreatePaymentInfo(db.Payment)
	if err != nil {
		t.Fatal("error creating payment info: " + err.Error())
	}

	paymnt := payment.Controller{Db: db, Validator: validatorRef, Logger: logger, ExtReq: request.ExternalRequest{
		Logger: logger,
		Test:   true,
	}}
	r := gin.Default()
	pvKey := utility.RandomString(20)
	pbKey := utility.RandomString(20)

	tests := []struct {
		Name         string
		RequestBody  interface{}
		ExpectedCode int
		Headers      map[string]string
		Message      string
		Reference    string
		Headless     bool
		SuccessPage  string
		FailurePage  string
		FundWallet   bool
	}{
		{
			Name:         "OK get status",
			ExpectedCode: http.StatusOK,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
			Message:     "Transaction payment successfully confirmed",
			Reference:   reference,
			Headless:    true,
			SuccessPage: "https://success.com",
			FailurePage: "https://failure.com",
			FundWallet:  true,
		},
		{
			Name:         "OK get status2",
			ExpectedCode: http.StatusOK,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
			Message:     "Transaction payment successfully confirmed",
			Reference:   reference2,
			Headless:    true,
			SuccessPage: "",
			FailurePage: "",
			FundWallet:  false,
		},
		{
			Name:         "no input",
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
		},
	}

	paymentApiUrl := r.Group(fmt.Sprintf("%v", "v2"), middleware.Authorize(db, paymnt.ExtReq, middleware.ApiType))
	{
		paymentApiUrl.GET("/pay/status", paymnt.GetStatus)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{
				Path:     "/v2/pay/status",
				RawQuery: fmt.Sprintf("reference=%v&headless=%v&success_page=%v&failure_page=%v&fund_wallet=%v", test.Reference, test.Headless, test.SuccessPage, test.FailurePage, test.FundWallet),
			}

			req, err := http.NewRequest(http.MethodGet, URI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			for i, v := range test.Headers {
				req.Header.Set(i, v)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tst.AssertStatusCode(t, rr.Code, test.ExpectedCode)

			data := tst.ParseResponse(rr)
			fmt.Println(data)

			code := int(data["code"].(float64))
			tst.AssertStatusCode(t, code, test.ExpectedCode)

			if test.Message != "" {
				message := data["message"]
				if message != nil {
					tst.AssertResponseMessage(t, message.(string), test.Message)
				} else {
					tst.AssertResponseMessage(t, "", test.Message)
				}

			}

		})

	}

}
func TestGetPaymentStatus(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := postgresql.Connection()
	var (
		muuid, _  = uuid.NewV4()
		accountID = uint(utility.GetRandomNumbersInRange(1000000000, 9999999999))
		testUser  = external_models.User{
			ID:           uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
			AccountID:    accountID,
			EmailAddress: fmt.Sprintf("testuser%v@qa.team", muuid.String()),
			PhoneNumber:  fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
			AccountType:  "individual",
			Firstname:    "test",
			Lastname:     "user",
			Username:     fmt.Sprintf("test_username%v", muuid.String()),
		}
	)

	auth_mocks.User = &testUser
	auth_mocks.BusinessProfile = &external_models.BusinessProfile{
		ID:        uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		AccountID: int(testUser.AccountID),
		Country:   "NG",
		Currency:  "NGN",
	}
	auth_mocks.ValidateAuthorizationRes = &external_models.ValidateAuthorizationDataModel{
		Status:  true,
		Message: "authorized",
		Data:    testUser,
	}
	auth_mocks.UserProfile = &external_models.UserProfile{
		ID:        uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		AccountID: int(testUser.AccountID),
		Country:   "NG",
		Currency:  "NGN",
	}

	auth_mocks.BusinessCharge = &external_models.BusinessCharge{
		ID:                  uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		BusinessId:          int(testUser.AccountID),
		Country:             "NG",
		Currency:            "NGN",
		BusinessCharge:      "0",
		VesicashCharge:      "2.5",
		ProcessingFee:       "0",
		PaymentGateway:      "rave",
		DisbursementGateway: "rave_momo",
		ProcessingFeeMode:   "fixed",
	}

	auth_mocks.Country = &external_models.Country{
		ID:           uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		Name:         "nigeria",
		CountryCode:  "NG",
		CurrencyCode: "NGN",
	}

	var (
		transactionID = utility.RandomString(20)
		partiesID     = utility.RandomString(20)
		milestoneID   = utility.RandomString(20)
		reference     = utility.RandomString(20)
		reference2    = utility.RandomString(20)
	)

	transactions_mocks.ListTransactionsByIDObj = &external_models.TransactionByID{
		ID:               uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		TransactionID:    transactionID,
		PartiesID:        partiesID,
		MilestoneID:      milestoneID,
		Title:            "title",
		Type:             "milestone",
		Description:      "description",
		Amount:           1000,
		Status:           "draft",
		Quantity:         1,
		InspectionPeriod: "2023-03-10",
		DueDate:          "1678924800",
		ShippingFee:      0,
		GracePeriod:      "1679097600",
		Currency:         "NGN",
		BusinessID:       int(testUser.AccountID),
		IsPaylinked:      false,
		Source:           "api",
		TransUssdCode:    74950,
		Recipients: []external_models.MileStoneRecipient{
			{
				AccountID:    9489042479,
				Amount:       500,
				EmailAddress: "test.qa.team",
				PhoneNumber:  "+23456789776789",
			},
		},
		EscrowWallet: "yes",
		Parties: map[string]external_models.TransactionParty{
			"buyer": {
				ID:                   uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
				TransactionPartiesID: partiesID,
				TransactionID:        transactionID,
				AccountID:            int(testUser.AccountID),
				Role:                 "buyer",
				Status:               "draft",
				RoleCapabilities: &map[string]interface{}{"approve": true,
					"can_receive":  false,
					"can_view":     true,
					"mark_as_done": false,
				},
			},
			"seller": {
				ID:                   uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
				TransactionPartiesID: partiesID,
				TransactionID:        transactionID,
				AccountID:            int(testUser.AccountID),
				Role:                 "seller",
				Status:               "draft",
				RoleCapabilities: &map[string]interface{}{"approve": true,
					"can_receive":  false,
					"can_view":     true,
					"mark_as_done": false,
				},
			},
		},
		Members: []external_models.PartyResponse{
			{
				PartyID:     utility.GetRandomNumbersInRange(1000000000, 9999999999),
				AccountID:   int(testUser.AccountID),
				AccountName: "name",
				Email:       "email@email.com",
				PhoneNumber: "+2340905964039",
				Role:        "buyer",
				Status:      "draft",
				AccessLevel: external_models.PartyAccessLevel{
					CanView:    true,
					CanReceive: false,
					MarkAsDone: false,
					Approve:    true,
				},
			},
			{
				PartyID:     utility.GetRandomNumbersInRange(1000000000, 9999999999),
				AccountID:   int(testUser.AccountID),
				AccountName: "name",
				Email:       "email@email.com",
				PhoneNumber: "+2340905964039",
				Role:        "seller",
				Status:      "draft",
				AccessLevel: external_models.PartyAccessLevel{
					CanView:    true,
					CanReceive: false,
					MarkAsDone: false,
					Approve:    true,
				},
			},
		},
		TotalAmount: 2000,
		Milestones: []external_models.MilestonesResponse{
			{
				Index:            1,
				MilestoneID:      utility.RandomString(20),
				Title:            "milestone title",
				Amount:           1000,
				Status:           "draft",
				InspectionPeriod: "2023-03-10",
				DueDate:          "1678924800",
				Recipients: []external_models.MilestonesRecipientResponse{
					{
						AccountID:   9489042479,
						AccountName: "name",
						Amount:      500,
						Email:       "test.qa.team",
						PhoneNumber: "+23456789776789",
					},
				},
			},
			{
				Index:            2,
				MilestoneID:      utility.RandomString(20),
				Title:            "milestone title 1",
				Amount:           1000,
				Status:           "draft",
				InspectionPeriod: "2023-03-10",
				DueDate:          "1678924800",
				Recipients: []external_models.MilestonesRecipientResponse{
					{
						AccountID:   9489042479,
						AccountName: "name",
						Amount:      500,
						Email:       "test.qa.team",
						PhoneNumber: "+23456789776789",
					},
				},
			},
		},
		IsDisputed: false,
	}

	paymentData := models.Payment{
		ID:               int64(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		TransactionID:    transactionID,
		PaymentID:        utility.RandomString(20),
		TotalAmount:      200,
		EscrowCharge:     0,
		IsPaid:           false,
		AccountID:        int64(accountID),
		BusinessID:       int64(accountID),
		Currency:         "NGN",
		DisburseCurrency: "NGN",
	}

	err := paymentData.CreatePayment(db.Payment)
	if err != nil {
		t.Fatal("errpr creating payment: " + err.Error())
	}
	paymentData2 := models.Payment{
		ID:               int64(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		TransactionID:    transactionID,
		PaymentID:        utility.RandomString(20),
		TotalAmount:      200,
		EscrowCharge:     0,
		IsPaid:           false,
		AccountID:        int64(accountID),
		BusinessID:       int64(accountID),
		Currency:         "NGN",
		DisburseCurrency: "NGN",
	}

	err = paymentData2.CreatePayment(db.Payment)
	if err != nil {
		t.Fatal("errpr creating payment: " + err.Error())
	}

	paymentCardInfo := models.PaymentCardInfo{
		AccountID:         int(testUser.AccountID),
		PaymentID:         paymentData.PaymentID,
		CcExpiryMonth:     "12",
		CcExpiryYear:      "2025",
		LastFourDigits:    "4566",
		Brand:             "VISA",
		IssuingCountry:    "nigeria",
		CardToken:         "",
		CardLifeTimeToken: "902be0d9be02be09b",
		Payload:           "",
	}
	err = paymentCardInfo.CreatePaymentCardInfo(db.Payment)
	if err != nil {
		t.Fatal("error creating payment card info: " + err.Error())
	}

	paymentAccount1 := models.PaymentAccount{
		PaymentAccountID: reference,
		TransactionID:    transactionID,
		PaymentID:        paymentData.PaymentID,
		IsUsed:           true,
		ExpiresAfter:     strconv.Itoa(int(time.Now().Add(72 * time.Hour).Unix())),
		BusinessID:       strconv.Itoa(int(testUser.AccountID)),
		AccountNumber:    config.GetConfig().Rave.MerchantId,
		BankCode:         "flutterwave",
		BankName:         "flutterwave",
		Status:           "ACTIVE",
	}
	err = paymentAccount1.CreatePaymentAccount(db.Payment)
	if err != nil {
		t.Fatal("error creating payment account: " + err.Error())
	}

	paymentAccount2 := models.PaymentAccount{
		PaymentID:        paymentData2.PaymentID,
		PaymentAccountID: reference2,
		IsUsed:           true,
		ExpiresAfter:     strconv.Itoa(int(time.Now().Add(72 * time.Hour).Unix())),
		BusinessID:       strconv.Itoa(int(testUser.AccountID)),
		AccountNumber:    "7727632865",
		BankCode:         "221",
		BankName:         "vesicash bank",
		Status:           "ACTIVE",
	}
	err = paymentAccount2.CreatePaymentAccount(db.Payment)
	if err != nil {
		t.Fatal("error creating payment account: " + err.Error())
	}

	paymentInfo1 := models.PaymentInfo{
		PaymentID:   paymentData.PaymentID,
		Reference:   reference,
		Status:      "pending",
		Gateway:     "rave",
		RedirectUrl: "",
	}

	err = paymentInfo1.CreatePaymentInfo(db.Payment)
	if err != nil {
		t.Fatal("error creating payment info: " + err.Error())
	}

	paymentInfo2 := models.PaymentInfo{
		PaymentID:   paymentData2.PaymentID,
		Reference:   reference2,
		Status:      "pending",
		Gateway:     "monnify",
		RedirectUrl: "",
	}

	err = paymentInfo2.CreatePaymentInfo(db.Payment)
	if err != nil {
		t.Fatal("error creating payment info: " + err.Error())
	}

	paymnt := payment.Controller{Db: db, Validator: validatorRef, Logger: logger, ExtReq: request.ExternalRequest{
		Logger: logger,
		Test:   true,
	}}
	r := gin.Default()
	pvKey := utility.RandomString(20)
	pbKey := utility.RandomString(20)

	tests := []struct {
		Name         string
		RequestBody  models.GetPaymentStatusRequest
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name:         "OK get new status",
			ExpectedCode: http.StatusOK,
			RequestBody: models.GetPaymentStatusRequest{
				Reference:  reference,
				Headless:   true,
				FundWallet: true,
			},
			Message: "Transaction payment successfully confirmed",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
		},
		{
			Name:         "OK get new status2",
			ExpectedCode: http.StatusOK,
			RequestBody: models.GetPaymentStatusRequest{
				Reference:  reference2,
				Headless:   true,
				FundWallet: false,
			},
			Message: "Transaction payment successfully confirmed",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
		},
		{
			Name:         "no input",
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
		},
	}

	paymentApiUrl := r.Group(fmt.Sprintf("%v", "v2"), middleware.Authorize(db, paymnt.ExtReq, middleware.ApiType))
	{
		paymentApiUrl.POST("/pay/new-status", paymnt.GetPaymentStatus)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/pay/new-status"}

			req, err := http.NewRequest(http.MethodPost, URI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			for i, v := range test.Headers {
				req.Header.Set(i, v)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tst.AssertStatusCode(t, rr.Code, test.ExpectedCode)

			data := tst.ParseResponse(rr)
			fmt.Println(data)

			code := int(data["code"].(float64))
			tst.AssertStatusCode(t, code, test.ExpectedCode)

			if test.Message != "" {
				message := data["message"]
				if message != nil {
					tst.AssertResponseMessage(t, message.(string), test.Message)
				} else {
					tst.AssertResponseMessage(t, "", test.Message)
				}

			}

		})

	}

}
