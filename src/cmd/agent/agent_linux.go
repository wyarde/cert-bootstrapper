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
	cert, err := getCert()
	if err != nil {
		return err
	}

	f, err := os.Open("/home")
	if err != nil {
		return err
	}
	defer f.Close()

	users, err := f.Readdirnames(0)
	if err != nil {
		return err
	}

	defaultPaths := []string{"/root", "/etc/skel"}
	paths := append(defaultPaths, users...)

	for _, path := range paths {
		log.WithField("path", path).Debug("  Adding Artifactory CLI configuration...")

		certPath := filepath.Join(path, ".jfrog/security/certs")
		err = os.MkdirAll(certPath, os.ModePerm)
		if err != nil {
			return err
		}

		f := filepath.Join(certPath, "cert.pem")
		err = os.WriteFile(f, cert, 0444)
		if err != nil {
			return err
		}
	}

	return nil
}
