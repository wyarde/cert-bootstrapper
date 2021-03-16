package main

import (
	_ "embed"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

//go:embed ssl/cacert.pem
var certBundle []byte

func configureNpm() error {
	log.Debug("Start of configureNpm")

	_ = os.Mkdir("/ssl", os.ModePerm)
	destinationFile := "/ssl/cacert.pem"
	err := os.WriteFile(destinationFile, certBundle, 0444)
	if err != nil {
		return err
	}

	log.Debug("End of configureNpm")
	return nil
}

func addCertToStore() error {
	log.Debug("Start of addCertToStore")

	cmd := exec.Command("certutil", "-addstore", "-f", "Root", "/cert.pem")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Debug("End of addCertToStore")

	return cmd.Run()
}

func bootstrap() error {
	var err error

	err = addCertToStore()
	if err != nil {
		return err
	}
	err = configureNpm()
	if err != nil {
		return err
	}
	return nil
}
