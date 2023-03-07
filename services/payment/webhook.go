package payment

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
)

func InitWebhook(extReq request.ExternalRequest, db postgresql.Databases, uri, event string, data map[string]interface{}, businessID int) error {
	if event == "" {
		event = "payment"
	}

	dataByte, err := json.Marshal(data)
	if err != nil {
		return err
	}

	webhook := models.Webhook{WebhookUri: strings.TrimSpace(uri), BusinessID: strconv.Itoa(businessID), Event: event, RequestPayload: string(dataByte)}
	err = webhook.CreateWebhook(db.Payment)
	if err != nil {
		return err
	}

	return nil
}

func fireWebhook(extReq request.ExternalRequest, db postgresql.Databases, webhook models.Webhook) error {
	var (
		method = "POST"
	)
	if extReq.Test {
		return nil
	}

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

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	extReq.Logger.Info("request for id: ", webhook.ID, webhook.WebhookUri, method, headers)

	res, err := client.Do(req)
	if err != nil {
		extReq.Logger.Error("client do error for id: ", webhook.ID, err.Error())
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		extReq.Logger.Error("reading body error for id: ", webhook.ID, err.Error())
		return err
	}

	webhook.RetryAt = time.Now()
	webhook.Tries = webhook.Tries + 1
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		webhook.IsReceived = true
		webhook.ResponsePayload = string(body)
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

	return nil
}

// TODO write cron job for webhooks
