package main

import (
	"github.com/wyarde/certificate-bootstrapper/cmd/agent/linux"
	"github.com/wyarde/certificate-bootstrapper/cmd/agent/windows"

	"os"
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

func main() {
	log.SetFormatter(&log.TextFormatter{TimestampFormat: time.RFC3339Nano})

	cert, err := os.ReadFile("/cert.pem")
	checkIfError(err)

	log.Info("Start of bootstrapper")
	switch runtime.GOOS {
	case "linux":
		err = linux.Bootstrap(cert)
	case "windows":
		err = windows.Bootstrap(cert)
	default:
		log.Error("Unknown operating system: ", runtime.GOOS)
		os.Exit(1)
	}

	checkIfError(err)
	log.Info("Done!")
}
