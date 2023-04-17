package cronjobs

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/payment-ms/services/payment"
)

func BankTransfer(extReq request.ExternalRequest, db postgresql.Databases) {
	pendingTransferFunding := models.PendingTransferFunding{Status: "pending"}
	transfers, err := pendingTransferFunding.GetAllBystatus(db.Payment)
	if err != nil {
		extReq.Logger.Error(fmt.Sprintf("error getting pending transfer funding %v", err.Error()))
		return
	}

	for _, item := range transfers {
		reference := item.Reference

		data, msg, code, err := payment.PaymentAccountMonnifyVerifyService(extReq, db, models.PaymentAccountMonnifyVerifyRequest{Reference: reference})
		if err != nil {
			extReq.Logger.Error("error cron job for bank transfer with reference: %v, data: %v, message:%v, code:%v, error:%v", reference, data, msg, code, err.Error())
			return
		}
		extReq.Logger.Info("cron job for bank transfer with reference: %v, data: %v, message:%v, code:%v", reference, data, msg, code)
	}

}
