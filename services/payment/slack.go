package payment

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/vesicash/payment-ms/internal/config"
)

func SlackNotify(channel, message string) error {
	api := slack.New(config.GetConfig().Slack.OauthToken)
	channelID := channel
	msg := message
	_, _, err := api.PostMessage(channelID, slack.MsgOptionText(msg, false))
	if err != nil {
		return fmt.Errorf("Error sending message: %v", err)
	}

	return nil
}
