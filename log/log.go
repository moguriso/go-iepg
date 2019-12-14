package log

import (
	"github.com/sirupsen/logrus"
)

type HDebugType int

var (
	Level logrus.Level = logrus.InfoLevel //default
	L                  = logrus.New()
)

func SetLogLevel(lv string) {
	tl, err := logrus.ParseLevel(lv)
	if err != nil {
		Level = logrus.DebugLevel
	} else {
		Level = tl
	}
}
