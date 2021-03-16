package main

import (
	"os"
	"runtime/debug"
	"time"

	log "github.com/sirupsen/logrus"
)

func checkIfError(err error) {
	if err == nil {
		return
	}

	log.Error(err)
	debug.PrintStack()

	os.Exit(1)
}

func main() {
	log.SetFormatter(&log.TextFormatter{TimestampFormat: time.RFC3339Nano})

	log.Info("Start of bootstrapper")

	err := bootstrap()
	checkIfError(err)

	log.Info("Done!")
}
