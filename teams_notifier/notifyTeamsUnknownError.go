package teams_notifier

import (
	"fmt"
	"os"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/messagecard"
)

func SendUnknownErrorToTeams(err error) {
	client := goteamsnotify.NewTeamsClient()
	webhook := os.Getenv("TEAMS_WEBHOOK_URL")

	card := messagecard.NewMessageCard()
	card.Title = "Unknown Error"
	card.Text = fmt.Sprintf("Script encountered an unknown error.<BR/>"+
		"**Error message**: %v", err)

	if err := client.Send(webhook, card); err != nil {
		fmt.Println("SendOrderErrorToTeams failed to send the error.")
	}

}
