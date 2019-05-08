package ecode

import (
	"fmt"
	log "go-open/library/log"
	"strconv"
	"sync"
)

var (
	// store error message. code as key, message as value
	_code = map[int]string{}
	mu    sync.RWMutex
)

type Code int

func (c Code) Error() string {
	return strconv.FormatInt(int64(c), 10)
}

func (c Code) Code() int { return int(c) }
func (c Code) Status() string {
	if c == OK {
		return "success"
	}
	return "failed"
}

func (c Code) Message() string {

	mu.RLock()
	defer mu.RUnlock()

	if msg, ok := _code[c.Code()]; ok {
		return msg
	}
	return c.Error()
}

func Add(code int, msg string) Code {
	if code < 0 {
		log.Info("ecode must greater than zero")
		panic("ecode must greater than zero")
	}

	mu.Lock()
	defer mu.Unlock()

	if _, ok := _code[code]; ok {
		log.Info(fmt.Sprintf("ecode: %d already exist", code))
		panic(fmt.Sprintf("ecode: %d already exist", code))
	}
	_code[code] = msg
	return Code(code)
}

func Cause(err error) Code {
	if err == nil {
		return OK
	}
	mu.RLock()
	defer mu.RUnlock()
	if c, ok := err.(Code); ok {
		return c
	}
	return UnknownCode
}
