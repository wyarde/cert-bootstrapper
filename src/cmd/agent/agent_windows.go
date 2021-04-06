package main

import (
	_ "embed"
	"os"
	"os/exec"
	"syscall"
)

func hideFile(filename string) {
	filenameW, err := syscall.UTF16PtrFromString(filename)
	if err != nil {
		syscall.SetFileAttributes(filenameW, syscall.FILE_ATTRIBUTE_HIDDEN)
	}
}

func addCertToStore() error {
	var err error

	cmd := exec.Command("certutil", "-addstore", "-f", "Root", "/.cert-bootstrapper/ssl/cert.pem")
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

	cmd := exec.Command("setx", "/m", "NODE_EXTRA_CA_CERTS", "/.cert-bootstrapper/ssl/cert.pem")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
