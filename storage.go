package main

import (
	"errors"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type ScheduledTip struct {
	Hour      int    `sql:"unique:hour_channel_id"`
	ChannelId string `sql:"unique:hour_channel_id"`
	Topic     string
}

type scheduledTipsStorage struct {
	db *pg.DB
}

func newScheduledTipsStorage(databaseUrl string) (*scheduledTipsStorage, error) {
	options, err := pg.ParseURL(databaseUrl)
	if err != nil {
		return nil, err
	}

	db := pg.Connect(options)
	if err := db.CreateTable(&ScheduledTip{}, &orm.CreateTableOptions{IfNotExists: true}); err != nil {
		return nil, err
	}

	return &scheduledTipsStorage{db}, nil
}

func (storage *scheduledTipsStorage) close() error {
	return storage.db.Close()
}

func (storage *scheduledTipsStorage) store(hour int, channelId string, topic string) error {
	if !isHourValid(hour) {
		return errors.New("invalid hour")
	}

	newScheduledTip := &ScheduledTip{Hour: hour, ChannelId: channelId, Topic: topic}
	_, err := storage.db.
		Model(newScheduledTip).
		OnConflict("(hour,channel_id) DO UPDATE").
		Set("topic = EXCLUDED.topic").
		Insert()

	return err
}

func (storage *scheduledTipsStorage) load(hour int) ([]ScheduledTip, error) {
	if !isHourValid(hour) {
		return nil, errors.New("invalid hour")
	}

	var res []ScheduledTip
	if err := storage.db.Model(&res).Where("hour = ?", hour).Select(); err != nil {
		return nil, err
	}

	return res, nil
}

func isHourValid(hour int) bool {
	return 0 <= hour && hour <= 23
}
