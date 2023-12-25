package logger

import "io"

type Logger interface {
	io.Closer
	Trace(msg string)
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Err(msg string)
	Fatal(msg string)
}

const (
	traceLogLevel = iota
	debugLogLevel
	infoLogLevel
	warnLogLevel
	errLogLevel
	fatalLogLevel
	noneLogLevel
)
