package main

import (
	"fmt"
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
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
		case "/list-topics":
			writeResponse(w, fmt.Sprintf("Avaliable topics are %v", strings.Join(listTopics(), ", ")))

		case "/get-tip":
			params := &slack.Msg{Text: slashCommand.Text}
			writeResponse(w, fmt.Sprintf("Your tip is: %s", getTip(params.Text)))

		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func writeResponse(w http.ResponseWriter, response string) {
	if _, err := w.Write([]byte(response)); err != nil {
		log.Warn("fail to write response")
	}
}
