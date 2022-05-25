package log

import (
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type LogWriter interface {
	Write(buf []byte) error
}

type StackTracer interface {
	StackTrace() errors.StackTrace
}

const (
	PanicLevel logrus.Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

func IsLevel(lvl logrus.Level) bool {
	return logrus.GetLevel() >= lvl
}

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}
func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	logrus.Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

func Tracef(format string, args ...interface{}) {
	logrus.Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func init() {
	lvl, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logrus.SetLevel(InfoLevel)
	} else {
		logrus.SetLevel(lvl)
	}
	logrus.SetFormatter(&prefixed.TextFormatter{})
}
