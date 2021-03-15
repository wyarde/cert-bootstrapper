package linux

import (
	"os"
	"os/exec"
)

func runUpdateCACertificates(cert []byte) error {
	destinationFile := "/usr/local/share/ca-certificates/cert.crt"

	err := os.WriteFile(destinationFile, cert, 0444)
	if err != nil {
		return err
	}

	cmd := exec.Command("update-ca-certificates", "--verbose")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func Bootstrap(cert []byte) error {
	return runUpdateCACertificates(cert)
}
