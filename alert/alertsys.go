package alert

import (
	"fmt"
	"log"

	"github.com/slack-go/slack"
)

// Alert system implemented using Slack
type AlertSys struct {
	channel string
	client  *slack.Client
}

// Create a new Slackbot alert system
func NewAlertSys(token, channel string) *AlertSys {
	return &AlertSys{
		channel: channel,
		client:  slack.New(token),
	}
}

// Send message to Slack channel
func (as *AlertSys) Notify(message string) error {
	log.Println("Sending message to Slack...")
	_, _, err := as.client.PostMessage(as.channel, slack.MsgOptionText(message, false))
	if err != nil {
		return fmt.Errorf("failed to send message to Slack: %w", err)
	}
	return nil
}
