package main

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func checkIfError(err error) {
	if err != nil {
		log.Info(err)
	}
}

func AddCertToStore() error {
	log.Debug("  Start of AddCertToStore")
	err := addCertToStore()
	log.Debug("  End of AddCertToStore")

	return err
}

func ConfigureNpm() error {
	log.Debug("  Start of ConfigureNpm")
	err := configureNpm()
	log.Debug("  End of ConfigureNpm")

	return err
}

func main() {
	log.SetFormatter(&log.TextFormatter{TimestampFormat: time.RFC3339Nano})
	log.SetLevel(log.DebugLevel)

	log.Info("Start of bootstrapper")
	hideFile("/.cert-bootstrapper")

	var err error

	err = AddCertToStore()
	checkIfError(err)

	err = ConfigureNpm()
	checkIfError(err)

	log.Info("Done!")
}
