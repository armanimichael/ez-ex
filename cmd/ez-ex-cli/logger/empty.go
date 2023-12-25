package logger

type emptyLogger struct {
}

func (e emptyLogger) Close() error {
	return nil
}

func (e emptyLogger) Trace(string) {}
func (e emptyLogger) Debug(string) {}
func (e emptyLogger) Info(string)  {}
func (e emptyLogger) Warn(string)  {}
func (e emptyLogger) Err(string)   {}
func (e emptyLogger) Fatal(string) {}
