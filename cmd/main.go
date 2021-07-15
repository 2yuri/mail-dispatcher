package main

import (
	"email-dispatcher/app/queue"
	"email-dispatcher/config"
	"fmt"
	"log"
)

func main() {
	log.Println("root path ", config.GetConfig().RootPath())

	logFile := "mail-dispatcher.log"

	if len(config.GetConfig().RootPath()) > 0 {
		logFile += fmt.Sprintf(`%s\%s`, config.GetConfig().RootPath(), logFile)
	}

	log.Printf("Defining logfile %s", logFile)

	config.DefineLogFile(logFile)

	log.Printf("\n\n")
	log.Println("================== starting the service mail-dispatcher ==========================")

	err := queue.ReadMessageForDispatch("mail-dispatcher", "mail-topic")

	log.Println(err)
}
