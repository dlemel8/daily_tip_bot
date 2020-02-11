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
	listTopicsCommand  = "/list-topics"
	getTipCommand      = "/get-tip"
	scheduleTipCommand = "/schedule-tip"
)

type slackBot struct {
	client  *slack.Client
	storage *scheduledTipsStorage
}

func newBot(botToken string, storage *scheduledTipsStorage) *slackBot {
	return &slackBot{slack.New(botToken), storage}
}

func (bot *slackBot) newSlashCommandHandler(signingSecret string) func(http.ResponseWriter, *http.Request) {
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
		case listTopicsCommand:
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
			hour, err := bot.convertToServerHour(slashCommand.UserID, hourStr)
			if err != nil {
				log.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			scheduledTip := &ScheduledTip{ChannelId: slashCommand.ChannelID, Topic: topic, Hour: hour}
			if err := bot.storage.store(scheduledTip); err != nil {
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

func (bot *slackBot) sendScheduledTips() error {
	hour, _, _ := time.Now().Clock()

	scheduledTips, err := bot.storage.loadByHour(hour)
	if err != nil {
		return err
	}

	log.Infof("found %d relevant scheduled tips", len(scheduledTips))
	lastHour := time.Now().Truncate(time.Hour)
	sent := 0

	for _, scheduledTip := range scheduledTips {
		if scheduledTip.Hour != hour {
			continue
		}

		if scheduledTip.LastSent.After(lastHour) {
			continue
		}

		message := fmt.Sprintf("Your %s daily tip is: %s", scheduledTip.Topic, getTip(scheduledTip.Topic))
		_, timestampStr, err := bot.client.PostMessage(scheduledTip.ChannelId, slack.MsgOptionText(message, false))
		if err != nil {
			log.Error(err)
			continue
		}
		sent += 1

		timestampInt, err := strconv.ParseFloat(timestampStr, 64)
		if err != nil {
			log.Error(err)
			continue
		}

		scheduledTip.LastSent = time.Unix(int64(timestampInt), 0)
		log.Infof("Message successfully sent at %s", scheduledTip.LastSent)
		err = bot.storage.store(&scheduledTip)
		if err != nil {
			log.Error(err)
			continue
		}
	}

	log.Infof("sent %d scheduled tips", sent)
	return nil
}

func (bot *slackBot) convertToServerHour(userId string, hourStr string) (int, error) {
	user, err := bot.client.GetUserInfo(userId)
	if err != nil {
		return 0, err
	}

	location, err := time.LoadLocation(user.TZ)
	if err != nil {
		return 0, err
	}

	userHourTime, err := time.ParseInLocation("15", hourStr, location)
	if err != nil {
		return 0, err
	}

	hour, err := strconv.Atoi(userHourTime.Local().Format("15"))
	if err != nil {
		return 0, err
	}
	return hour, nil
}

func writeResponse(w http.ResponseWriter, response string) {
	if _, err := w.Write([]byte(response)); err != nil {
		log.Warn("fail to write response")
	}
}
