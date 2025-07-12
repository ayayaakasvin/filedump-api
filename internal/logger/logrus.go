package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func SetupLogger() *logrus.Logger {
	logFile, err := os.OpenFile(".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	seperator := fmt.Sprintf("=== Logging started at %s ===\n", time.Now().Format("2006-01-02 15:04:05"))
	logFile.Write([]byte(seperator))

	multiwriter := io.MultiWriter(os.Stdout, logFile)

	logger := logrus.New()
	logger.SetOutput(multiwriter)
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.SetLevel(logrus.InfoLevel)

	logger.Info("Logger has been set up")
	return logger
}
