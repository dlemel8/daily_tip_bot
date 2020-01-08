package main

import (
	"errors"
	"fmt"
	"github.com/prometheus/common/log"
	"os"

	"github.com/nlopes/slack"
)

func run() error {
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		return errors.New("please set SLACK_TOKEN env var")
	}

	api := slack.New(token)
	if _, err := api.AuthTest(); err != nil {
		return fmt.Errorf("token is not valid: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
