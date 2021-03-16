package main

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func runUpdateCACertificates(cert []byte) error {
	log.Debug("Start of runUpdateCACertificates")

	destinationFile := "/usr/local/share/ca-certificates/cert.crt"
	err := os.WriteFile(destinationFile, cert, 0444)
	if err != nil {
		return err
	}

	cmd := exec.Command("update-ca-certificates", "--verbose")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Debug("End of runUpdateCACertificates")
	return cmd.Run()
}

func bootstrap() error {
	cert, err := os.ReadFile("/cert.pem")
	if err != nil {
		return err
	}

	return runUpdateCACertificates(cert)
}
