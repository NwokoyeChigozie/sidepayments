package payment

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/config"
)

func SlackNotify(extReq request.ExternalRequest, channel, message string) error {
	api := slack.New(config.GetConfig().Slack.OauthToken)
	channelID := channel
	msg := message
	_, _, err := api.PostMessage(channelID, slack.MsgOptionText(msg, false))
	if err != nil {
		extReq.Logger.Error("error sending message to slack", err.Error())
		return fmt.Errorf("error sending message: %v", err)
	}

	return nil
}
