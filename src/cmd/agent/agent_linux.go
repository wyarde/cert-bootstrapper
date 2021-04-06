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

	destinationFile := "/usr/local/share/ca-certificates/cert.crt"
	err = os.WriteFile(destinationFile, cert, 0444)
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
			path := prefix + "/etc"
			npmrcFile := path + "/npmrc"

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

func configureArtifactoryCli() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	certFolder := filepath.Join(home, ".jfrog/security/certs")
	err = os.MkdirAll(certFolder, os.ModePerm)
	if err != nil {
		return err
	}

	destinationFile := filepath.Join(certFolder, "cert.pem")

	cert, err := getCert()
	if err != nil {
		return err
	}

	err = os.WriteFile(destinationFile, cert, 0444)
	if err != nil {
		return err
	}

	return nil
}
