package main

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func run() error {
	port := os.Getenv("PORT")
	if port == "" {
		return errors.New("please set PORT env var")
	}

	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")
	if signingSecret == "" {
		return errors.New("please set SLACK_SIGNING_SECRET env var")
	}

	http.HandleFunc("/slack", newSlashCommandHandler(signingSecret))

	log.Info("Server listening")
	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
