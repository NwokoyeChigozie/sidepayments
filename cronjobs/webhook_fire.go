package cronjobs

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/payment-ms/services/payment"
)

func WebhookFire(extReq request.ExternalRequest, db postgresql.Databases) {
	var (
		webhook = models.Webhook{IsAbandoned: false, IsReceived: false}
	)

	webhooks, err := webhook.GetAllByIsAbandonedAndIsReceived(db.Payment)
	if err != nil {
		extReq.Logger.Error(fmt.Sprintf("error getting webhoks error: %v", err.Error()))
		return
	}

	for _, item := range webhooks {
		payment.FireWebhook(extReq, db, item)
	}
}
