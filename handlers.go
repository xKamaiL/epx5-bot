package main

import (
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
)

func handlers(s *session.Session) {
	s.AddHandler(func(i *gateway.MessageCreateEvent) {})
}
