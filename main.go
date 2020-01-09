package main

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func run() error {
	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")
	if signingSecret == "" {
		return errors.New("please set SLACK_SIGNING_SECRET env var")
	}

	http.HandleFunc("/slack", newSlashCommandHandler(signingSecret))

	log.Info("Server listening")
	return http.ListenAndServe(":8080", nil)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
