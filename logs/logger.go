package logs

import (
	"log"
)

// LogLevel to indicates log should output or not
type LogLevel int

const (
	// InfoLog for info logs
	InfoLog = iota
	// WarningLog for warning logs
	WarningLog
	// ErrorLog for error logs
	ErrorLog
	// FatalLog for fatal logs
	FatalLog
)

var logLevel LogLevel

// SetLogLevel set the global log level, log entry with a less log level won't print to output
func SetLogLevel(ll LogLevel) {
	logLevel = ll
}

// Logger used in cas library
type Logger interface {
	// Infof output to INFO logs with format and arguments
	Infof(format string, args ...interface{})

	// Info output to INFO logs
	Info(v ...interface{})

	// Infoln output to INFO logs with a new line
	Infoln(v ...interface{})

	// Warningf output WARNING logs
	Warningf(format string, args ...interface{})

	// Warning output to WARNING logs
	Warning(v ...interface{})

	// Warningln output to WARNING logs with a new line
	Warningln(v ...interface{})

	// Errorf output ERROR logs
	Errorf(format string, args ...interface{})

	// Error output to ERROR logs
	Error(v ...interface{})

	// Errorln output to ERROR logs with a new line
	Errorln(v ...interface{})

	// Fatalf output FATAL logs, then exist
	Fatalf(format string, args ...interface{})

	// Fatal output to FATAL logs
	Fatal(v ...interface{})

	// Fatalln output to FATAL logs with a new line, then exist
	Fatalln(v ...interface{})
}

var defaultLogger Logger = &stdLogger{}

// SetLogger can change the logger implementation
func SetLogger(lg Logger) {
	defaultLogger = lg
}

type stdLogger struct {
}

func (*stdLogger) Infof(format string, args ...interface{}) {
	if logLevel > InfoLog {
		return
	}

	log.Printf(format, args...)
}

func (l *stdLogger) Info(v ...interface{}) {
	if logLevel > InfoLog {
		return
	}

	log.Print(v...)
}

func (l *stdLogger) Infoln(v ...interface{}) {
	if logLevel > InfoLog {
		return
	}

	log.Println(v...)
}

func (*stdLogger) Warningf(format string, args ...interface{}) {
	if logLevel > WarningLog {
		return
	}

	log.Printf(format, args...)
}

func (l *stdLogger) Warning(v ...interface{}) {
	if logLevel > WarningLog {
		return
	}

	log.Print(v...)
}

func (l *stdLogger) Warningln(v ...interface{}) {
	if logLevel > WarningLog {
		return
	}

	log.Println(v...)
}

func (*stdLogger) Errorf(format string, args ...interface{}) {
	if logLevel > ErrorLog {
		return
	}

	log.Printf(format, args...)
}

func (l *stdLogger) Error(v ...interface{}) {
	if logLevel > ErrorLog {
		return
	}

	log.Print(v...)
}

func (l *stdLogger) Errorln(v ...interface{}) {
	if logLevel > ErrorLog {
		return
	}

	log.Println(v...)
}

func (*stdLogger) Fatalf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l *stdLogger) Fatal(v ...interface{}) {
	log.Fatal(v...)
}

func (l *stdLogger) Fatalln(v ...interface{}) {
	log.Fatalln(v...)
}

var _ Logger = (*stdLogger)(nil)

// Infof output to INFO logs with format and arguments
func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

// Info output to INFO logs
func Info(v ...interface{}) {
	defaultLogger.Info(v...)
}

// Infoln output to INFO logs with a new line
func Infoln(v ...interface{}) {
	defaultLogger.Infoln(v...)
}

// Warningf output WARNING logs
func Warningf(format string, args ...interface{}) {
	defaultLogger.Warningf(format, args...)
}

// Warning output to WARNING logs
func Warning(v ...interface{}) {
	defaultLogger.Warning(v...)
}

// Warningln output to WARNING logs with a new line
func Warningln(v ...interface{}) {
	defaultLogger.Warningln(v...)
}

// Errorf output ERROR logs
func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

// Error output to ERROR logs
func Error(v ...interface{}) {
	defaultLogger.Error(v...)
}

// Errorln output to ERROR logs with a new line
func Errorln(v ...interface{}) {
	defaultLogger.Errorln(v...)
}

// Fatalf output FATAL logs, then exist
func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}

// Fatal output to FATAL logs
func Fatal(v ...interface{}) {
	defaultLogger.Fatal(v...)
}

// Fatalln output to FATAL logs with a new line, then exist
func Fatalln(v ...interface{}) {
	defaultLogger.Fatalln(v...)
}
