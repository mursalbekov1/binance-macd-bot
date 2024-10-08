package logging

import (
	"github.com/google/uuid"
	"log"
	"os"
	"time"
)

const (
	YYYYMMDD  = "2006-01-02"
	HHMMSS24h = "15:04:05"
)

func CurrentDatetime() string {
	return time.Now().Format(YYYYMMDD + "-" + HHMMSS24h)
}

func CustomLog(prefix string, uid string) (*log.Logger, *os.File) {
	logFile := InitLogFile()
	logger := log.New(logFile, CurrentDatetime()+`: `+prefix+` uid `+uid+`:`, log.Lshortfile)

	return logger, logFile
}

func InitLogFile() *os.File {
	LOG_FILE := "/var/log/binanceapp.log"
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
