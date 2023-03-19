package utils

import (
	"log"
	"os"
	"time"
)

func GetHttpLogFile() *log.Logger {
	t := time.Now()
	startTime := t.Format("2006-01-02 15:04:05")
	logFile, err := os.Create("httpServerLog/ErrorFile " + startTime + ".log")
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(logFile, "", log.LstdFlags)

	return logger
}
