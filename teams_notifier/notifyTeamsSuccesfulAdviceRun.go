package teams_notifier

import (
	"fmt"
	"os"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/messagecard"
)

type MagasinAdviceInfo struct {
	AdviceNumber int
	HouseNumber  string
}

func SendAdviceSuccesToTeams(Advices []MagasinAdviceInfo, AdviceType string) error {
	client := goteamsnotify.NewTeamsClient()
	webhook := os.Getenv("TEAMS_WEBHOOK_URL")

	card := messagecard.NewMessageCard()
	if len(Advices) == 0 {
		fmt.Println("No advices to send to teams?")
	}

	card.Title = AdviceType + ": Succesful Advice run"

	cardText := "Script has run and sent the following advices to magasin.<BR/>"
	for _, advice := range Advices {
		cardText += fmt.Sprintf("**Advice**: %v - House: %v <BR/>", advice.AdviceNumber, advice.HouseNumber)
	}

	card.Text = cardText

	if err := client.Send(webhook, card); err != nil {
		return fmt.Errorf("SendValidationErrorToTeams failed to send the error. Error: %v", err)
	}
	return nil
}
