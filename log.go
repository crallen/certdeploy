package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

type logFormatter struct{}

func (f *logFormatter) Format(entry *log.Entry) ([]byte, error) {
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
	log.SetOutput(os.Stdout)
	log.SetFormatter(new(logFormatter))
}
