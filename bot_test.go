package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type dummySocketmodeClient struct{}

func (d dummySocketmodeClient) Ack(_ socketmode.Request, _ ...interface{}) {}
func (d dummySocketmodeClient) Run() error                                 { return nil }

type dummySlackClient struct {
	err  bool
	body string
}

func (d *dummySlackClient) PostMessage(_ string, opts ...slack.MsgOption) (_, _ string, err error) {
	if d.err {
		err = fmt.Errorf("an error")
	}

	_, r, _ := slack.UnsafeApplyMsgOptions("", "", "", opts...)

	d.body = r.Get("text")

	return
}

type dummyHTTPClient struct {
	status int
	err    bool
}

func (d dummyHTTPClient) Do(_ *http.Request) (resp *http.Response, err error) {
	if d.err {
		err = fmt.Errorf("an error")
	}

	resp = new(http.Response)
	resp.StatusCode = d.status
	resp.Body = io.NopCloser(bytes.NewBufferString(`{"message":{"result":{"translatedText":"foo"}}}}`))

	return
}

func TestNew(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			t.Fatal(err)
		}
	}()

	New("", "", "", "")
}

func TestBot_Process(t *testing.T) {
	for _, test := range []struct {
		name        string
		message     string
		slackClient slackClient
		httpclient  httpClient
		expect      string
	}{
		{"happy path, english message", "Good morning!", &dummySlackClient{}, dummyHTTPClient{}, "foo"},
		{"happy path, korean message", "오.. 신기합니다.", &dummySlackClient{}, dummyHTTPClient{}, "foo"},
		{"papago errors", "message", &dummySlackClient{}, dummyHTTPClient{err: true}, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			b := Bot{
				s:         dummySocketmodeClient{},
				slack:     test.slackClient,
				client:    test.httpclient,
				eventChan: make(chan socketmode.Event),
			}

			go func() {
				b.eventChan <- socketmode.Event{
					Type: socketmode.EventTypeEventsAPI,
					Data: slackevents.EventsAPIEvent{
						Type: slackevents.CallbackEvent,
						InnerEvent: slackevents.EventsAPIInnerEvent{
							Data: &slackevents.MessageEvent{
								Text:            test.message,
								TimeStamp:       "0",
								ThreadTimeStamp: "0",
							},
						},
					},
					Request: &socketmode.Request{},
				}

				close(b.eventChan)
			}()

			b.Process()

			body := b.slack.(*dummySlackClient).body
			if body != test.expect {
				t.Errorf("expected %q, received %q", test.expect, body)
			}
		})
	}

}
