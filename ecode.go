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

// A Ecode is an unsigned 32-bit error code as defined in the gRPC spec.
type Ecode int32

type Codes interface {
	Err() error
	Code() Ecode
	Message() string
	Details() []interface{}
}

func Register(m map[int]string) {
	_messages.Store(m)
}

func (e Ecode) Error() string {
	return strconv.FormatInt(int64(e), 10)
}

// Code return error code
func (e Ecode) Code() int { return int(e) }

// Message return error message
func (e Ecode) Message() string {
	if cm, ok := _messages.Load().(map[int]string); ok {
		if msg, ok := cm[e.Code()]; ok {
			return msg
		}
	}
	return e.Error()
}

// Details return details.
func (e Ecode) Details() []interface{} { return nil }

func New(e int) Ecode {
	if e <= 1000 {
		panic("business ecode must greater than 1000")
	}
	return add(e)
}

func add(e int) Ecode {
	if _, ok := _codes[e]; ok {
		panic(fmt.Sprintf("ecode: %d already exist", e))
	}
	_codes[e] = struct{}{}
	return Ecode(e)
}

// Cause cause from error to Codes.
func Cause(err error) Codes {
	if err == nil {
		return Convert(nil)
	}

	return Convert(err)
}

// EqualError equal error
func EqualError(e Ecode, err error) bool {
	return Cause(err).Code() == e
}
