package token

import (
	xtime "go-open/library/time"
)

type Token struct {
}

type TokenConfig struct {
	Secret     string
	Expiration xtime.Duration
}
