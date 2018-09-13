package errors

import (
	"bytes"
	"runtime/debug"
)

const (
	InternalZenatonError = "InternalZenatonError"
	ExternalZenatonError = "ExternalZenatonError"
	ScheduledBoxError    = "ScheduledBoxError"
)

type ZenatonError interface {
	Error() string
	Trace() string
	Name() string
}

type zenatonErrorImp struct {
	name    string
	trace   string
	message string
}

func (ze *zenatonErrorImp) Error() string {
	return ze.message
}

func (ze *zenatonErrorImp) Trace() string {
	return ze.trace
}

func (ze *zenatonErrorImp) Name() string {
	return ze.name
}

func New(name, message string) ZenatonError {
	return NewWithOffset(name, message, 4)
}

func Wrap(name string, err error) ZenatonError {
	return WrapWithOffset(name, err, 4)
}

func NewWithOffset(name, message string, offset int) ZenatonError {
	trace := getTraceWithOffset(offset)
	return &zenatonErrorImp{
		name:    name,
		message: message,
		trace:   trace,
	}
}

func WrapWithOffset(name string, err error, offset int) ZenatonError {

	if err == nil {
		return nil
	}

	trace := getTraceWithOffset(offset)
	return &zenatonErrorImp{
		name:    name,
		message: err.Error(),
		trace:   trace,
	}
}

func getTraceWithOffset(offset int) string {
	stack := debug.Stack()
	parts := bytes.Split(stack, []byte("\n"))
	partsMinusRemoved := parts[offset*2+1:]
	stack = bytes.Join(partsMinusRemoved, []byte("\n"))
	return string(stack)
}
