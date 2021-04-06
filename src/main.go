package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/ahmetalpbalkan/dlog"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
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

func execInContainer(ctx context.Context, cli *client.Client, id string, commands []string) (output string, err error) {
	execID, err := cli.ContainerExecCreate(ctx, id, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          commands,
		Privileged:   true,
	})

	if err != nil {
		return "", err
	}

	res, err := cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return "", err
	}

	r := dlog.NewReader(res.Reader)
	s := bufio.NewScanner(r)
	for s.Scan() {
		output += fmt.Sprintf("  > %s\n", s.Text())
	}

	return output, nil
}

// File holds the filename and its content
type File struct {
	Name    string
	Content []byte
	Mode    int64
}

func generateArchive(files ...File) io.Reader {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Size: int64(len(file.Content)),
			Mode: file.Mode,
		}

		err := tw.WriteHeader(hdr)
		checkIfError(err)

		_, err = tw.Write(file.Content)
		checkIfError(err)
	}

	err := tw.Close()
	checkIfError(err)

	return buf
}

//go:embed bin/agent
var agent []byte

// OsData type describes Osspecific data requiredfor bootstrapping
type OsData struct {
	Archive io.Reader
	Command []string
}

func getOsData(containerOs string, cert []byte) (osData OsData, err error) {
	certFile := File{
		Name:    "/.cert-bootstrapper/ssl/cert.pem",
		Content: cert,
		Mode:    444,
	}

	agentFileName := "/.cert-bootstrapper/bootstrap-agent"
	if containerOs == "windows" {
		agentFileName = fmt.Sprintf("%s.exe", agentFileName)
	}

	agentFile := File{
		Name:    agentFileName,
		Content: agent,
		Mode:    555,
	}

	osData = OsData{
		Archive: generateArchive(agentFile, certFile),
		Command: []string{agentFileName},
	}

	return osData, nil
}

func isHyperVContainer(ctx context.Context, cli *client.Client, id string) bool {
	info, err := cli.ContainerInspect(ctx, id)
	checkIfError(err)

	return info.ContainerJSONBase.HostConfig.Isolation == "hyperv"
}

func getContainerOs(ctx context.Context, cli *client.Client, id string) string {
	info, err := cli.ContainerInspect(ctx, id)
	checkIfError(err)

	return info.ContainerJSONBase.Platform
}

func bootstrap(ctx context.Context, cli *client.Client, id string, from string, cert []byte) {
	log.WithFields(log.Fields{
		"id":   id,
		"from": from,
	}).Info("Bootstrapping container...")

	hostname, err := os.Hostname()
	checkIfError(err)

	log.WithField("hostname", hostname).Info()
	if id == hostname {
		log.Info("This is me, don't want to touch myself. No action taken.")
		return
	}

	if isHyperVContainer(ctx, cli, id) {
		log.Warning("This container is running in Hyper-V isolation, which is not supported. No action taken.")
		return
	}

	containerOs := getContainerOs(ctx, cli, id)
	log.WithField("os", containerOs).Info()

	if containerOs != runtime.GOOS {
		log.Warningf("Container OS %s doesn't match host OS %s. No action taken.", containerOs, runtime.GOOS)
		return
	}

	osData, err := getOsData(containerOs, cert)
	checkIfError(err)

	log.Info("Copying files into container...")
	err = cli.CopyToContainer(ctx, id, "/", osData.Archive, types.CopyToContainerOptions{})
	checkIfError(err)

	log.Info("Running bootstrap script...")
	output, err := execInContainer(ctx, cli, id, osData.Command)
	checkIfError(err)

	fmt.Println("==== Output bootstrap script ====")
	fmt.Println(output)
	fmt.Println("=================================")
	log.Info("Done!")
	fmt.Println("\n")
}

func main() {
	log.SetFormatter(&log.TextFormatter{TimestampFormat: time.RFC3339Nano})

	certFilename := "ssl/cert.pem"
	if len(os.Args) >= 2 {
		certFilename = os.Args[1]
	}

	log.WithField("certFilename", certFilename).Info("Reading certificate file...")

	cert, err := ioutil.ReadFile(certFilename)
	checkIfError(err)

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	checkIfError(err)

	ctx := context.Background()
	checkIfError(err)

	msgs, errs := cli.Events(ctx, types.EventsOptions{})

	log.Info("Listening for new containers...")

	for {
		select {
		case err := <-errs:
			log.Error(err)
			os.Exit(1)
		case msg := <-msgs:
			if msg.Status == "start" {
				bootstrap(ctx, cli, msg.ID[0:12], msg.From, cert)
			}
		}
	}
}
