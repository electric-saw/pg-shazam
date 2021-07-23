package config

import (
	"time"
)

type Health struct {
	Timeout string `yaml:"timeout"`
	Retries int    `yaml:"retries"`
}

func NewHealth() *Health {
	return &Health{
		Retries: 3,
		Timeout: "20s",
	}
}

func (h *Health) TimeoutDuration() (time.Duration, error) {
	return time.ParseDuration(h.Timeout)
}
