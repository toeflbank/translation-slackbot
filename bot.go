package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/abadojack/whatlanggo"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type Response struct {
	Message struct {
		Result struct {
			Text string `json:"translatedText"`
		}
	}
}

type socketmodeClient interface {
	Ack(socketmode.Request, ...interface{})
	Run() error
}

type slackClient interface {
	PostMessage(string, ...slack.MsgOption) (string, string, error)
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Bot struct {
	eventChan chan socketmode.Event
	s         socketmodeClient
	slack     slackClient

	client       httpClient
	clientID     string
	clientSecret string
}

func New(slackBotToken, slackAppToken, clientID, clientSecret string) (b Bot, err error) {
	b.slack = slack.New(slackBotToken,
		slack.OptionAppLevelToken(slackAppToken),
	)

	b.s = socketmode.New(
		b.slack.(*slack.Client),
	)

	b.eventChan = b.s.(*socketmode.Client).Events

	go b.s.Run()

	b.client = http.DefaultClient

	b.clientID = clientID
	b.clientSecret = clientSecret

	return
}

// Process will:
//  1. Listen to slack events
//  2. On channel message, detect character encoding
//  3. Perform appropriate translation
//  4. Respond as thread reply
func (b Bot) Process() error {
	for evt := range b.eventChan {
		switch evt.Type {
		case socketmode.EventTypeEventsAPI:
			eventsAPIEvent, _ := evt.Data.(slackevents.EventsAPIEvent)

			switch eventsAPIEvent.Type {
			case slackevents.CallbackEvent:
				innerEvent := eventsAPIEvent.InnerEvent

				switch ev := innerEvent.Data.(type) {
				case *slackevents.MessageEvent:
					if ev.BotID != "" {
						continue
					}

					var (
						body string
						err  error
					)

					if whatlanggo.Detect(ev.Text).Lang == whatlanggo.Kor {
						body, err = b.toEN(ev.Text)
					} else {
						body, err = b.toKO(ev.Text)
					}

					if err != nil {
						log.Print(err)

						continue
					}

					ts := ev.TimeStamp
					if ev.ThreadTimeStamp != "" {
						ts = ev.ThreadTimeStamp
					}

					b.slack.PostMessage(ev.Channel,
						slack.MsgOptionText(body, false),
						slack.MsgOptionTS(ts),
					)
				}
			}

			b.s.Ack(*evt.Request)

		}
	}

	return nil
}

func (b Bot) toKO(msg string) (string, error) {
	return b.translate("en", "ko", msg)
}

func (b Bot) toEN(msg string) (string, error) {
	return b.translate("ko", "en", msg)
}

func (b Bot) translate(from, to, msg string) (s string, err error) {
	form := url.Values{}
	form.Set("source", from)
	form.Set("target", to)
	form.Set("text", msg)

	req, err := http.NewRequest("POST", "https://openapi.naver.com/v1/papago/n2mt", strings.NewReader(form.Encode()))
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("X-Naver-Client-Id", b.clientID)
	req.Header.Add("X-Naver-Client-Secret", b.clientSecret)

	resp, err := b.client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	r := new(Response)
	err = dec.Decode(r)

	if err != nil {
		return
	}

	s = r.Message.Result.Text

	return
}
