package logging

import (
	"github.com/google/uuid"
	"log"
	"os"
)

func CustomLog(prefix string, uid string) (*log.Logger, *os.File) {
	logFile := InitLogFile()
	logger := log.New(logFile, prefix+`,uid `+uid+`:`, log.Lshortfile)

	return logger, logFile
}

func InitLogFile() *os.File {
	LOG_FILE := "/var/logs/binanceapp.log"
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panicln(err)
	}
	return logFile
}

func GenerateString(size int) string {
	if size != 32 {
		return uuid.New().String()
	}
	return uuid.New().String()
}
