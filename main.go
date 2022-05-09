package main

import (
	"context"
	"log"
	"os"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
)

var (
	token string
)

func main() {
	token = os.Getenv("BOT_TOKEN")
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	s := session.New("Bot " + token)
	s.AddHandler(func(c *gateway.MessageCreateEvent) {
		log.Println(c.Author.Username, "sent", c.Content)
	})
	// Add the needed Gateway intents.
	s.AddIntents(gateway.IntentGuildMessages)
	s.AddIntents(gateway.IntentDirectMessages)

	if err := s.Open(context.Background()); err != nil {
		return err
	}
	defer s.Close()

	u, err := s.Me()
	if err != nil {
		return err
	}

	log.Println("Started as", u.Username)
	// Block forever.
	select {}
}
