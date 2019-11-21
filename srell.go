package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hekonsek/osexit"
	"github.com/nlopes/slack"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func main() {
	NewSrell().Connect()
}

type Srell struct {
	pwd string
}

func NewSrell() *Srell {
	return &Srell{pwd:"/"}
}

func (srell *Srell) Connect() {
	slackToken := os.Getenv("SLACK_TOKEN")
	if slackToken == "" {
		osexit.ExitBecauseError("SLACK_TOKEN cannot be empty.")
	}

	api := slack.New(slackToken, slack.OptionDebug(true))
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})
	go func() {
		osexit.ExitOnError(r.Run())
	}()

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
						cmd.Dir = srell.pwd
						if x == "cd" {
							srell.pwd = commandParts[2]
							status := fmt.Sprintf("`STATUS: OK (%d)`", osexit.UnixExitCodeOK)
							rtm.SendMessage(&slack.OutgoingMessage{Type: "message", Channel: m["channel"].(string), Text: status})
						} else {
							out, err := cmd.CombinedOutput()
							if err != nil {
								println(err.Error())
								status := ""
								if exiterr, ok := err.(*exec.ExitError); ok {
									if s, ok := exiterr.Sys().(syscall.WaitStatus); ok {
										status = fmt.Sprintf("`STATUS: (%d)`", s.ExitStatus())
									}
								} else {
									status = fmt.Sprintf("`STATUS: Error executing command (%s)`", err.Error())
								}
								rtm.SendMessage(&slack.OutgoingMessage{Type: "message", Channel: m["channel"].(string), Text: "```" + string(out) + "```\n" + status})
							} else {
								status := fmt.Sprintf("`STATUS: OK (%d)`", osexit.UnixExitCodeOK)
								rtm.SendMessage(&slack.OutgoingMessage{Type: "message", Channel: m["channel"].(string), Text: "```" + string(out) + "```\n" + status})
							}
						}
					}
				}
			}
		}
	}
}
