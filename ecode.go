package ecode

import (
	"fmt"
	"strconv"
	"sync/atomic"
)

var (
	_messages atomic.Value         // NOTE: stored map[int]string
	_codes    = map[int]struct{}{} // register codes.
)

type Codes interface {
	Err() error
	Code() Code
	Message() string
	Details() []interface{}
}

func (e Code) Error() string {
	return strconv.FormatInt(int64(e), 10)
}

// Code return error code
func (e Code) Code() int { return int(e) }

// Message return error message
func (e Code) Message() string {
	if cm, ok := _messages.Load().(map[int]string); ok {
		if msg, ok := cm[e.Code()]; ok {
			return msg
		}
	}
	return e.Error()
}

func Register(m map[int]string) {
	_messages.Store(m)
}

// Details return details.
func (e Code) Details() []interface{} { return nil }

func New(e int) Code {
	if e <= 1000 {
		panic("business ecode must greater than 1000")
	}
	return add(e)
}

func add(e int) Code {
	if _, ok := _codes[e]; ok {
		panic(fmt.Sprintf("ecode: %d already exist", e))
	}
	_codes[e] = struct{}{}
	return Code(e)
}

// Cause cause from error to Codes.
func Cause(err error) Codes {
	if err == nil {
		return Convert(nil)
	}

	return Convert(err)
}

// EqualError equal error
func EqualError(e Code, err error) bool {
	return Cause(err).Code() == e
}
