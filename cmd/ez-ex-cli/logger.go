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
	Info(format string)
	Warn(format string)
	Err(format string)
	Debug(format string)
	Trace(format string)
	Fatal(format string)
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

func (f FileLogger) Info(format string) {
	f.print("[INFO]", format, f.flag)
}

func (f FileLogger) Warn(format string) {
	f.print("[WARNING]", format, f.flag)
}

func (f FileLogger) Err(format string) {
	f.print("[ERROR]", format, f.flag)
}

func (f FileLogger) Debug(format string) {
	f.print("[DEBUG]", format, f.flag|log.Llongfile)
}

func (f FileLogger) Trace(format string) {
	f.print("[TRACE]", format, f.flag|log.Llongfile)
}

func (f FileLogger) Fatal(format string) {
	f.print("[FATAL]", format, f.flag)
	os.Exit(1)
}

func (f FileLogger) print(prefix string, format string, flag int) {
	f.logger.SetFlags(flag)
	f.logger.SetPrefix(prefix + " ")
	f.logger.Println(format)
}
