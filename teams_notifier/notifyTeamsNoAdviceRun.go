package teams_notifier

import (
	"fmt"
	"os"
	"time"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/messagecard"
)

func SendNoAdviceToTeams(AdviceType string) error {
	client := goteamsnotify.NewTeamsClient()
	webhook := os.Getenv("TEAMS_WEBHOOK_URL")

	card := messagecard.NewMessageCard()
	card.Title = AdviceType + ": No Advice run"
	card.Text = fmt.Sprintf("Script has run at %v but found no new advices.<BR/>", time.Now().Format("2006-01-02 15:04:05"))

	if err := client.Send(webhook, card); err != nil {
		return fmt.Errorf("SendValidationErrorToTeams failed to send the error. Error: %v", err)
	}
	return nil
}
