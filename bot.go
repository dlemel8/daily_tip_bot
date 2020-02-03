package main

import (
	"fmt"
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	listTopicCommand   = "/list-topics"
	getTipCommand      = "/get-tip"
	scheduleTipCommand = "/schedule-tip"
)

func newSlashCommandHandler(signingSecret string, storage *scheduledTipsStorage) func(http.ResponseWriter, *http.Request) {
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
		case listTopicCommand:
			writeResponse(w, fmt.Sprintf("Avaliable topics are %v", strings.Join(listTopics(), ", ")))

		case getTipCommand:
			writeResponse(w, fmt.Sprintf("Your tip is: %s", getTip(slashCommand.Text)))

		case scheduleTipCommand:
			params := strings.Split(slashCommand.Text, " ")
			if len(params) != 2 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			hourStr, topic := params[0], params[1]
			hour, err := strconv.Atoi(hourStr)
			if err != nil {
				log.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if err := storage.store(hour, slashCommand.ChannelID, topic); err != nil {
				log.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			writeResponse(w, fmt.Sprintf("Schedule a new tip from topic %s on hour %d ", topic, hour))

		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

func sendScheduledTips(botToken string, storage *scheduledTipsStorage) error {
	api := slack.New(botToken)
	hour, _, _ := time.Now().Clock()

	scheduledTips, err := storage.load(hour)
	if err != nil {
		return err
	}

	log.Infof("found %d relevant scheduled tips", len(scheduledTips))

	for _, scheduledTip := range scheduledTips {
		if scheduledTip.Hour != hour {
			continue
		}

		message := fmt.Sprintf("Your %s daily tip is: %s", scheduledTip.Topic, getTip(scheduledTip.Topic))
		if _, timestamp, err := api.PostMessage(
			scheduledTip.ChannelId,
			slack.MsgOptionText(message, false)); err != nil {
			return err
		} else {
			log.Info("Message successfully sent to at %s", timestamp)
		}
	}

	return nil
}

func writeResponse(w http.ResponseWriter, response string) {
	if _, err := w.Write([]byte(response)); err != nil {
		log.Warn("fail to write response")
	}
}
