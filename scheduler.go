package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type ScheduledTip struct {
	Hour      int    `json:"hour"`
	ChannelId string `json:"channel_id"`
	Topic     string `json:"topic"`
}

const path = "local.json"

var empty = ScheduledTip{}

type scheduledTipsStorage struct {
	content map[int]ScheduledTip
}

func (storage *scheduledTipsStorage) store(hour int, channelId string, topic string) error {
	if err := storage.validation(hour); err != nil {
		return err
	}

	newScheduledTip := ScheduledTip{hour, channelId, topic}
	if currentScheduledTip, ok := storage.content[hour]; ok {
		if newScheduledTip == currentScheduledTip {
			return nil
		}
	}

	storage.content[hour] = newScheduledTip
	marshalled, err := json.Marshal(storage.content)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, marshalled, 0600)
}

func (storage *scheduledTipsStorage) load(hour int) (ScheduledTip, error) {
	if err := storage.validation(hour); err != nil {
		return empty, err
	}

	if len(storage.content) == 0 {
		marshalled, err := ioutil.ReadFile(path)
		if err != nil {
			return empty, err
		}
		storage.content = make(map[int]ScheduledTip)
		if err := json.Unmarshal(marshalled, &storage.content); err != nil {
			return empty, err
		}
	}

	if currentScheduledTip, ok := storage.content[hour]; ok {
		return currentScheduledTip, nil
	}
	return empty, nil
}

func (storage *scheduledTipsStorage) validation(hour int) error {
	if hour < 0 || hour > 23 {
		return errors.New("invalid hour")
	}

	if storage.content == nil {
		storage.content = make(map[int]ScheduledTip)
	}

	return nil
}
