package main

import (
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func hideFile(filename string) {
}

func getCert() ([]byte, error) {
	return os.ReadFile("/.cert-bootstrapper/ssl/cert.pem")
}

func addCertToStore() error {
	var err error

	cert, err := getCert()
	if err != nil {
		return err
	}

	f := "/usr/local/share/ca-certificates/cert.crt"
	err = os.WriteFile(f, cert, 0444)
	if err != nil {
		return err
	}

	cmd := exec.Command("update-ca-certificates")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func configureNpm() error {
	var err error

	cmd := exec.Command("npm", "config", "set", "cafile", "/.cert-bootstrapper/ssl/cert.pem")
	err = cmd.Run()
	if err != nil {
		// Since npm isn't there, generate config in common default locations in case npm is added later
		log.Debug("  Npm NOT installed. Creating configuration files at the default locations.")

		for _, prefix := range []string{"/usr", "/usr/local"} {
			path := filepath.Join(prefix, "etc")
			npmrcFile := filepath.Join(path, "npmrc")

			log.WithField("npmrcFile", npmrcFile).Debug()

			_ = os.Mkdir(path, os.ModeDir)
			f, err := os.OpenFile(npmrcFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = f.WriteString("cafile=/.cert-bootstrapper/ssl/cert.pem\n")
			if err != nil {
				return err
			}
		}
	}

	return nil
}
