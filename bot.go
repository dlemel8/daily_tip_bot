package main

import (
	"fmt"
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

func newSlashCommandHandler(signingSecret string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
		slashCommand, err := slack.SlashCommandParse(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = verifier.Ensure(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		switch slashCommand.Command {
		case "/get-tip":
			params := &slack.Msg{Text: slashCommand.Text}
			response := fmt.Sprintf("You asked for a tip from source %v", params.Text)
			if _, err := w.Write([]byte(response)); err != nil {
				log.Warn("fail to write response")
			}

		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
