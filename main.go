package main

import (
	"os"
)

const (
	Nick = "my-bot"
	Chan = "#dashboard"
)

var (
	Username  = os.Getenv("SASL_USER")
	Password  = os.Getenv("SASL_PASSWORD")
	Server    = os.Getenv("SERVER")
	VerifyTLS = os.Getenv("VERIFY_TLS") == "true"
)

func main() {
	c, err := New(Username, Password, Server, VerifyTLS)
	if err != nil {
		panic(err)
	}

	panic(c.bottom.Client.Connect())
}
