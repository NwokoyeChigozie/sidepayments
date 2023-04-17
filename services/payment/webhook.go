package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/payment-ms/utility"
)

func InitWebhook(extReq request.ExternalRequest, db postgresql.Databases, uri, event string, data map[string]interface{}, businessID int) error {
	if event == "" {
		event = "payment"
	}

	dataByte, err := json.Marshal(data)
	if err != nil {
		return err
	}

	webhook := models.Webhook{WebhookUri: strings.TrimSpace(uri), BusinessID: strconv.Itoa(businessID), Event: event, RequestPayload: string(dataByte), IsAbandoned: false, IsReceived: false}
	err = webhook.CreateWebhook(db.Payment)
	if err != nil {
		return err
	}

	FireWebhook(extReq, db, webhook)
	return nil
}

func FireWebhook(extReq request.ExternalRequest, db postgresql.Databases, webhook models.Webhook) error {
	var (
		method = "POST"
	)

	if extReq.Test {
		return nil
	}

	if webhook.RetryAt != "" {
		retryAtUnix, _ := strconv.Atoi(webhook.RetryAt)
		retryAt := time.Unix(int64(retryAtUnix), 0)
		if retryAt.After(time.Now()) {
			return nil
		}
	}
	businessId, _ := strconv.Atoi(webhook.BusinessID)
	apiKey, _ := GetAccessTokenByBusinessID(extReq, businessId)

	data := map[string]interface{}{}

	json.Unmarshal([]byte(webhook.RequestPayload), &data)

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(data)
	if err != nil {
		extReq.Logger.Error("webhook error for id: ", webhook.ID, err.Error())
		return err
	}

	extReq.Logger.Info("request for id: ", webhook.WebhookUri, data, webhook.RequestPayload)

	client := &http.Client{}
	req, err := http.NewRequest(method, webhook.WebhookUri, buf)
	if err != nil {
		extReq.Logger.Error("request creation error for id: ", webhook.ID, err.Error())
		return err
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	if apiKey.PrivateKey != "" {
		headerSignature := fmt.Sprintf("%v:%v", apiKey.PrivateKey, businessId)
		hmacResult := utility.Sha256Hmac(apiKey.PrivateKey, []byte(headerSignature))
		headers["X-Vesicash-Webhook-Secret"] = hmacResult
		headers["User-Agent"] = "Vesicash Agent/2.0"
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	extReq.Logger.Info("request for id: ", webhook.ID, webhook.WebhookUri, method, headers)

	res, err := client.Do(req)
	if err != nil {
		extReq.Logger.Error("client do error for id: ", webhook.ID, err.Error())
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		extReq.Logger.Error("reading body error for id: ", webhook.ID, err.Error())
		return err
	}

	webhook.RetryAt = strconv.Itoa(int(time.Now().Add(24 * time.Hour).Unix()))
	webhook.Tries = webhook.Tries + 1
	webhook.ResponseCode = strconv.Itoa(res.StatusCode)
	webhook.ResponsePayload = string(body)

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		webhook.IsReceived = true
	} else {

		webhook.IsReceived = false
	}

	defer res.Body.Close()

	if webhook.Tries > 10 {
		webhook.IsAbandoned = true
	}

	err = webhook.UpdateAllFields(db.Payment)
	if err != nil {
		return err
	}

	extReq.Logger.Info(fmt.Sprintf("Webhook #%v FIRED, RESPONSE CODE: %v", webhook.ID, res.StatusCode))

	return nil
}

// TODO write cron job for webhooks
