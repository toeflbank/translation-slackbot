package main

import (
	"os"
)

var (
	slackBotToken = os.Getenv("SLACK_BOT_TOKEN")
	slackAppToken = os.Getenv("SLACK_APP_TOKEN")

	naverClientID     = os.Getenv("NAVER_CLIENT_ID")
	naverClientSecret = os.Getenv("NAVER_CLIENT_SECRET")
)

func main() {
	c, err := New(slackBotToken, slackAppToken, naverClientID, naverClientSecret)
	if err != nil {
		panic(err)
	}

	panic(c.Process())
}
