package main

import (
	ezex "github.com/armanimichael/ez-ex"
	"io"
	"log"
	"os"
	"path"
)

type Logger interface {
	io.Closer
	Info(msg string)
	Warn(msg string)
	Err(msg string)
	Debug(msg string)
	Trace(msg string)
	Fatal(msg string)
}

type FileLogger struct {
	logger *log.Logger
	file   *os.File
	flag   int
}

func NewFileLogger() Logger {
	home, _ := os.UserHomeDir()
	logsFile, err := os.OpenFile(path.Join(home, ezex.UserDataDir, "logs"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Fatalf("Error opening log file: %s", err)
	}

	defaultFlag := log.LstdFlags

	return FileLogger{
		logger: log.New(logsFile, "", defaultFlag),
		file:   logsFile,
		flag:   defaultFlag,
	}
}

func (f FileLogger) Close() error {
	return f.file.Close()
}

func (f FileLogger) Info(msg string) {
	f.print("[INFO]", msg, f.flag)
}

func (f FileLogger) Warn(msg string) {
	f.print("[WARNING]", msg, f.flag)
}

func (f FileLogger) Err(msg string) {
	f.print("[ERROR]", msg, f.flag)
}

func (f FileLogger) Debug(msg string) {
	f.print("[DEBUG]", msg, f.flag|log.Llongfile)
}

func (f FileLogger) Trace(msg string) {
	f.print("[TRACE]", msg, f.flag|log.Llongfile)
}

func (f FileLogger) Fatal(msg string) {
	f.print("[FATAL]", msg, f.flag)
}

func (f FileLogger) print(prefix string, msg string, flag int) {
	f.logger.SetFlags(flag)
	f.logger.SetPrefix(prefix + " ")
	f.logger.Println(msg)
}
