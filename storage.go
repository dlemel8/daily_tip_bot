package main

import (
	"errors"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"time"
)

type ScheduledTip struct {
	ChannelId string `sql:"unique:channel_topic"`
	Topic     string `sql:"unique:channel_topic"`
	Hour      int
	LastSent  time.Time
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

func (storage *scheduledTipsStorage) store(scheduledTip *ScheduledTip) error {
	if !isHourValid(scheduledTip.Hour) {
		return errors.New("invalid hour")
	}

	_, err := storage.db.
		Model(scheduledTip).
		OnConflict("(channel_id,topic) DO UPDATE").
		Set("hour = EXCLUDED.hour, last_sent = EXCLUDED.last_sent").
		Insert()

	return err
}

func (storage *scheduledTipsStorage) loadByHour(hour int) ([]ScheduledTip, error) {
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
