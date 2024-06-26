package main

import (
	"flag"
	"log"

	tgClient "github.com/fumkaa/read-adviser-bot/clients/telegram"
	event_consumer "github.com/fumkaa/read-adviser-bot/consumer/event-consumer"
	"github.com/fumkaa/read-adviser-bot/events/telegram"
	"github.com/fumkaa/read-adviser-bot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {
	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("servise is stopped: %w", err)
	}
}

func mustToken() string {
	token := flag.String("tg-bot-token",
		"",
		"token for access to telegram bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
