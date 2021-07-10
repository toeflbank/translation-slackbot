package main

import (
	"github.com/jspc-bots/bottom"
	"github.com/lrstanley/girc"
)

type Bot struct {
	bottom bottom.Bottom
}

func New(user, password, server string, verify bool) (b Bot, err error) {
	b.bottom, err = bottom.New(user, password, server, verify)
	if err != nil {
		return
	}

	b.bottom.Client.Handlers.Add(girc.CONNECTED, func(c *girc.Client, e girc.Event) {
		c.Cmd.Join(Chan)
	})

	router := bottom.NewRouter()
	router.AddRoute(`hello\s+world\!`, b.hello)

	b.bottom.Middlewares.Push(router)

	return
}

// hello will respond to a sender on the appropriate channel (an irc channel if the
// message was recieved in a channel, and a message if received directly
func (b Bot) hello(sender, channel string, groups []string) error {
	b.bottom.Client.Cmd.Messagef(channel, "And hello to you too %s", sender)

	return nil
}
