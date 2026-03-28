package config

import "time"

type Worker struct {
	TimeInterval time.Duration `split_words:"true" default:"1s"`
	EventCount   uint16        `split_words:"true" default:"50"`
	// ThreadCount is the number of concurrent workers processing events.
	ThreadCount     uint8         `split_words:"true" default:"3"`
}
