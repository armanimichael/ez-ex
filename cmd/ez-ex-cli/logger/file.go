package logger

import (
	ezex "github.com/armanimichael/ez-ex"
	"log"
	"os"
	"path"
)

// FileLogger logs to ezex.UserDataDir/logs
type FileLogger struct {
	logger *log.Logger
	level  int
	file   *os.File
	flag   int
}

// NewFileLogger instantiates a new file logger
// note: level must be between 0 (trace) and 6 (none)
//
//	0 = trace
//	1 = debug
//	2 = info
//	3 = warn
//	4 = err
//	5 = fatal
//	6 = none
func NewFileLogger(level int) Logger {
	if level > 6 || level < 0 {
		panic("log level must be between 0 (trace) and 6 (none)")
	}

	if level == 6 {
		return emptyLogger{}
	}

	home, _ := os.UserHomeDir()
	logsFile, err := os.OpenFile(path.Join(home, ezex.UserDataDir, "logs"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Fatalf("Error opening log file: %s", err)
	}

	defaultFlag := log.LstdFlags

	return FileLogger{
		logger: log.New(logsFile, "", defaultFlag),
		level:  level,
		file:   logsFile,
		flag:   defaultFlag,
	}
}

func (f FileLogger) Close() error {
	return f.file.Close()
}

func (f FileLogger) Info(msg string) {
	if f.level <= infoLogLevel {
		f.print("[INFO]", msg, f.flag)
	}
}

func (f FileLogger) Warn(msg string) {
	if f.level <= warnLogLevel {
		f.print("[WARNING]", msg, f.flag)
	}
}

func (f FileLogger) Err(msg string) {
	if f.level <= errLogLevel {
		f.print("[ERROR]", msg, f.flag)
	}
}

func (f FileLogger) Debug(msg string) {
	if f.level <= debugLogLevel {
		f.print("[DEBUG]", msg, f.flag|log.Llongfile)
	}
}

func (f FileLogger) Trace(msg string) {
	if f.level <= traceLogLevel {
		f.print("[TRACE]", msg, f.flag|log.Llongfile)
	}
}

func (f FileLogger) Fatal(msg string) {
	if f.level <= fatalLogLevel {
		f.print("[FATAL]", msg, f.flag)
	}
}

func (f FileLogger) print(prefix string, msg string, flag int) {
	f.logger.SetFlags(flag)
	f.logger.SetPrefix(prefix + " ")
	f.logger.Println(msg)
}
