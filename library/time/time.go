package xtime

import (
	"time"
)

type Duration time.Duration

func (d *Duration) UnmarshalText(text []byte) error {
	var (
		du  time.Duration
		err error
	)
	if du, err = time.ParseDuration(string(text)); err == nil {
		*d = Duration(du)
	}
	return err
}
