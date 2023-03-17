package teams_notifier

import (
	"fmt"
	"os"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/messagecard"
)

func SendRequestsReturnErrorToTeams(requestName string, requestType string, response string, responseBody string, api string) {
	client := goteamsnotify.NewTeamsClient()
	webhook := os.Getenv("TEAMS_WEBHOOK_URL")

	card := messagecard.NewMessageCard()
	card.Title = "Request Error"
	card.Text = fmt.Sprintf("**API:** %s <BR/>"+
		"**Request Type:** %s<BR/>"+
		"**Request Name:** %s <BR/>"+
		"**Response **: %s <BR/>"+
		"**ResponseBody **: %s <BR/>", api, requestName, requestType, response, responseBody)

	if err := client.Send(webhook, card); err != nil {
		fmt.Println("SendVatCodeErrorToTeams failed to send the error.")
	}
}
