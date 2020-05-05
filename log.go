package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type logFormatter struct{}

func (f *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	s := entry.Message
	if ns, ok := entry.Data["namespace"]; ok {
		s = fmt.Sprintf("[%v] %s", ns, s)
	}
	if cluster, ok := entry.Data["cluster"]; ok {
		s = fmt.Sprintf("[%v] %s", cluster, s)
	}
	return append([]byte(s), '\n'), nil
}

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(new(logFormatter))
}
