package test_payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/external/mocks/auth_mocks"
	"github.com/vesicash/payment-ms/external/mocks/transactions_mocks"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/controller/payment"
	"github.com/vesicash/payment-ms/pkg/middleware"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	tst "github.com/vesicash/payment-ms/tests"
	"github.com/vesicash/payment-ms/utility"
)

func TestCreatePayment(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := postgresql.Connection()
	var (
		muuid, _  = uuid.NewV4()
		token, _  = uuid.NewV4()
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

	paymnt := payment.Controller{Db: db, Validator: validatorRef, Logger: logger, ExtReq: request.ExternalRequest{
		Logger: logger,
		Test:   true,
	}}
	r := gin.Default()

	tests := []struct {
		Name         string
		RequestBody  models.CreatePaymentRequest
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK create payment",
			RequestBody: models.CreatePaymentRequest{
				TransactionID: transactionID,
				TotalAmount:   2000,
				ShippingFee:   0,
				EscrowCharge:  10,
				Currency:      "NGN",
			},
			ExpectedCode: http.StatusCreated,
			Message:      "Created",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
		{
			Name:         "no input",
			RequestBody:  models.CreatePaymentRequest{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
	}

	paymentAuthUrl := r.Group(fmt.Sprintf("%v", "v2"), middleware.Authorize(db, paymnt.ExtReq, middleware.AuthType))
	{
		paymentAuthUrl.POST("/create", paymnt.CreatePayment)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/create"}

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

func TestCreatePaymentHeadless(t *testing.T) {
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

	paymnt := payment.Controller{Db: db, Validator: validatorRef, Logger: logger, ExtReq: request.ExternalRequest{
		Logger: logger,
		Test:   true,
	}}
	r := gin.Default()
	pvKey := utility.RandomString(20)
	pbKey := utility.RandomString(20)

	tests := []struct {
		Name         string
		RequestBody  models.CreatePaymentHeadlessRequest
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK create payment",
			RequestBody: models.CreatePaymentHeadlessRequest{
				AccountID:    int(testUser.AccountID),
				TotalAmount:  2000,
				EscrowCharge: 10,
				Currency:     "NGN",
			},
			ExpectedCode: http.StatusCreated,
			Message:      "Created",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
		}, {
			Name: "no account id",
			RequestBody: models.CreatePaymentHeadlessRequest{
				TotalAmount:  2000,
				EscrowCharge: 10,
				Currency:     "NGN",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
		}, {
			Name: "no amount",
			RequestBody: models.CreatePaymentHeadlessRequest{
				AccountID:    int(testUser.AccountID),
				EscrowCharge: 10,
				Currency:     "NGN",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
		},
		{
			Name:         "no input",
			RequestBody:  models.CreatePaymentHeadlessRequest{},
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
		paymentApiUrl.POST("/create/headless", paymnt.CreatePaymentHeadless)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/create/headless"}

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

func TestEditPayment(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := postgresql.Connection()
	var (
		muuid, _  = uuid.NewV4()
		token, _  = uuid.NewV4()
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
		PaymentID:        utility.RandomString(20),
		TransactionID:    transactionID,
		TotalAmount:      2000,
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

	paymnt := payment.Controller{Db: db, Validator: validatorRef, Logger: logger, ExtReq: request.ExternalRequest{
		Logger: logger,
		Test:   true,
	}}
	r := gin.Default()

	tests := []struct {
		Name         string
		RequestBody  models.EditPaymentRequest
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK edit payment",
			RequestBody: models.EditPaymentRequest{
				PaymentID:    paymentData.PaymentID,
				EscrowCharge: 20,
			},
			ExpectedCode: http.StatusOK,
			Message:      "Edited",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		}, {
			Name: "no paymentid",
			RequestBody: models.EditPaymentRequest{
				EscrowCharge: 20,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		}, {
			Name: "OK escrow charge",
			RequestBody: models.EditPaymentRequest{
				EscrowCharge: 20,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
		{
			Name:         "no input",
			RequestBody:  models.EditPaymentRequest{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
	}

	paymentAuthUrl := r.Group(fmt.Sprintf("%v", "v2"), middleware.Authorize(db, paymnt.ExtReq, middleware.AuthType))
	{
		paymentAuthUrl.PATCH("/edit", paymnt.EditPayment)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/edit"}

			req, err := http.NewRequest(http.MethodPatch, URI.String(), &b)
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

func TestVerifyTransactionPayment(t *testing.T) {
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
		PaymentID:        utility.RandomString(20),
		TransactionID:    transactionID,
		TotalAmount:      2000,
		EscrowCharge:     0,
		IsPaid:           true,
		AccountID:        int64(accountID),
		BusinessID:       int64(accountID),
		Currency:         "NGN",
		DisburseCurrency: "NGN",
		PaymentMadeAt:    time.Now(),
	}

	err := paymentData.CreatePayment(db.Payment)
	if err != nil {
		t.Fatal("errpr creating payment: " + err.Error())
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
		RequestBody  models.VerifyTransactionPaymentRequest
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK verify transaction payment",
			RequestBody: models.VerifyTransactionPaymentRequest{
				TransactionID: transactionID,
			},
			ExpectedCode: http.StatusOK,
			Message:      "Payment Details Retrieved",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
		},
		{
			Name:         "no input",
			RequestBody:  models.VerifyTransactionPaymentRequest{},
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
		paymentApiUrl.PATCH("/verify", paymnt.VerifyTransactionPayment)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/verify"}

			req, err := http.NewRequest(http.MethodPatch, URI.String(), &b)
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

func TestDeletePayment(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := postgresql.Connection()
	var (
		muuid, _  = uuid.NewV4()
		token, _  = uuid.NewV4()
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
		PaymentID:        utility.RandomString(20),
		TransactionID:    transactionID,
		TotalAmount:      2000,
		EscrowCharge:     0,
		IsPaid:           true,
		AccountID:        int64(accountID),
		BusinessID:       int64(accountID),
		Currency:         "NGN",
		DisburseCurrency: "NGN",
		PaymentMadeAt:    time.Now(),
	}

	err := paymentData.CreatePayment(db.Payment)
	if err != nil {
		t.Fatal("errpr creating payment: " + err.Error())
	}

	paymnt := payment.Controller{Db: db, Validator: validatorRef, Logger: logger, ExtReq: request.ExternalRequest{
		Logger: logger,
		Test:   true,
	}}
	r := gin.Default()

	tests := []struct {
		Name         string
		RequestBody  interface{}
		ExpectedCode int
		Headers      map[string]string
		Message      string
		PaymentID    string
	}{
		{
			Name:         "OK delete payment",
			ExpectedCode: http.StatusOK,
			Message:      "Deleted",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
			PaymentID: paymentData.PaymentID,
		},
	}

	paymentAuthUrl := r.Group(fmt.Sprintf("%v", "v2"), middleware.Authorize(db, paymnt.ExtReq, middleware.AuthType))
	{
		paymentAuthUrl.DELETE("/delete/:id", paymnt.DeletePayment)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/delete/" + test.PaymentID}

			req, err := http.NewRequest(http.MethodDelete, URI.String(), &b)
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
