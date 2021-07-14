package main

import (
	"os"
)

var (
	SlackBotToken = os.Getenv("SLACK_BOT_TOKEN")
	SlackAppToken = os.Getenv("SLACK_APP_TOKEN")

	NaverClientId     = os.Getenv("NAVER_CLIENT_ID")
	NaverClientSecret = os.Getenv("NAVER_CLIENT_SECRET")
)

func main() {
	c, err := New(SlackBotToken, SlackAppToken, NaverClientId, NaverClientSecret)
	if err != nil {
		panic(err)
	}

	panic(c.Process())
}
