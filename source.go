package main

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type tipSource interface {
	randomTip() string
}

type inMemorySource struct {
	tips []string
}

func (s *inMemorySource) randomTip() string {
	if s.tips == nil {
		return ""
	}
	return s.tips[rand.Intn(len(s.tips))]
}
