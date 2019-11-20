package main

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"os"
	"os/exec"
	"strings"
)

func main() {
	BotConnect()
}

func BotConnect() {
	slackToken := os.Getenv("SLACK_TOKEN")
	if slackToken == "" {
		println("SLACK_TOKEN cannot be empty.")
		os.Exit(1)
	}

	api := slack.New(slackToken, slack.OptionDebug(true))
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch msg.Data.(type) {
		case *slack.UnmarshallingErrorEvent:
			eventText := fmt.Sprintf("%+v\n", msg)
			eventPayloadWithClosingBracketParts := strings.Split(eventText, "message\":")
			if len(eventPayloadWithClosingBracketParts) != 2 {
				println("Skipping....")
			} else {
				eventPayloadWithClosingBracket := eventPayloadWithClosingBracketParts[1]
				eventPayload := eventPayloadWithClosingBracket[0 : len(eventPayloadWithClosingBracket)-2]
				var m map[string]interface{}
				err := json.Unmarshal([]byte(eventPayload), &m)
				if err != nil {
					println(err.Error())
				} else {
					if strings.HasPrefix(m["text"].(string), "shell") {
						commandParts := strings.Split(m["text"].(string), " ")
						x := commandParts[1]
						y := commandParts[2:]
						cmd := exec.Command(x, y...)
						if x == "cd" {
							cmd.Dir = commandParts[2]
						}
						out, err := cmd.CombinedOutput()
						if err != nil {
							println(err.Error())
						} else {
							rtm.SendMessage(&slack.OutgoingMessage{Type: "message", Channel: m["channel"].(string), Text: "```" + string(out) + "```"})
						}
					}
				}
			}
		}
	}
}
