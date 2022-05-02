package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Logger struct {
	*logrus.Logger
	f *os.File
}

func NewLogger(path string) Logger {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	l.SetFormatter(&logrus.JSONFormatter{})
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		l.SetOutput(os.Stdout)
	}
	l.SetOutput(f)
	return Logger{l, f}
}

func (l Logger) CloseFile() {
	if l.f != nil {
		l.f.Close()
	}
}
