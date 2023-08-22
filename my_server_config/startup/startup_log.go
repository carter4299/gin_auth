package startup

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func s_log() *logrus.Logger {
	log := logrus.New()

	log.AddHook(&contextHook{})
	log.SetFormatter(&cformatter{})
	file, err := os.OpenFile("./my_server_config/startup/startup.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(io.MultiWriter(file))
	log.SetLevel(logrus.DebugLevel)

	return log
}

type contextHook struct{}

func (hook *contextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *contextHook) Fire(entry *logrus.Entry) error {
	if _, file, line, ok := runtime.Caller(5); ok {
		entry.Data["file"] = fmt.Sprintf("%s:%d", file, line)
	}
	return nil
}

type cformatter struct{}

func (f *cformatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Timestamp - LogLevel - Message
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	log := fmt.Sprintf("%s - %s - %s\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(log), nil
}
