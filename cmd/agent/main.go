package main

import (
	"os"
	"os/exec"
	"runtime"
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

func runUpdateCACertificates() {
	input, err := os.ReadFile("/cert.pem")
	checkIfError(err)

	destinationFile := "/usr/local/share/ca-certificates/cert.crt"
	err = os.WriteFile(destinationFile, input, 0444)
	if err != nil {
		log.Warning("Error creating file", destinationFile)
		return
	}

	cmd := exec.Command("update-ca-certificates", "--verbose")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	checkIfError(err)
}

func main() {
	log.SetFormatter(&log.TextFormatter{TimestampFormat: time.RFC3339Nano})

	log.Info("Start of bootstrapper")
	switch runtime.GOOS {
	case "linux":
		runUpdateCACertificates()
	case "windows":
		log.Info("Hello World")
	default:
		log.Error("Unknown operating system: ", runtime.GOOS)
		os.Exit(1)
	}

	log.Info("Done!")
}
