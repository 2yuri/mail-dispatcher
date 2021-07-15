package config

import (
	"log"
	"os"
)

func DefineLogFile(logFile string) *os.File {

	fileLog, err := os.OpenFile(logFile, os.O_RDWR | os.O_CREATE | os.O_APPEND , 0666)

	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(fileLog)
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	return fileLog
}
