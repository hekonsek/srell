package main

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"os/exec"
	"strings"
)

func main() {
	ListenAndResponse()
}

func ListenAndResponse() {
	api := slack.New("xoxb-826173663410-840819711239-vFyTDYDV592MZjZc2ANVDnRr", slack.OptionDebug(true))
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
							out, err := cmd.CombinedOutput()
							if err != nil {
								println(err.Error())
							} else {
								rtm.SendMessage(&slack.OutgoingMessage{Type: "message", Channel: m["channel"].(string), Text: "```" + string(out) + "```" })
							}
						}
					}
				}
		}
	}
}