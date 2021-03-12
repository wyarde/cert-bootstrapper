package main

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime/debug"

	"github.com/ahmetalpbalkan/dlog"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
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
		Privileged:   true,
		Cmd:          commands,
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

//go:embed scripts/bootstrap.sh
var bootstrapLinux string

//go:embed scripts/bootstrap.ps1
var bootstrapWindows string

// OsData type describes Osspecific data requiredfor bootstrapping
type OsData struct {
	Archive io.Reader
	Command []string
}

func generateArchive(input ...string) io.Reader {
	out, err := archive.Generate(input...)
	checkIfError(err)

	return out
}

func getOsData(cert []byte) (archives map[string]OsData, err error) {

	osData := map[string]OsData{
		"linux": {
			Archive: generateArchive("bootstrap.sh", bootstrapLinux, "cert.pem", string(cert)),
			Command: []string{"sh", "/bootstrap.sh"},
		},
		"windows": {
			Archive: generateArchive("bootstrap.ps1", bootstrapWindows, "cert.pem", string(cert)),
			Command: []string{"powershell", "-NoProfile", "-Command", "/bootstrap.ps1"},
		},
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

func bootstrap(ctx context.Context, cli *client.Client, id string, from string, osData map[string]OsData) {
	log.WithFields(log.Fields{
		"id":   id[0:12],
		"from": from,
	}).Info("Bootstrapping container...")

	if isHyperVContainer(ctx, cli, id) {
		log.Warning("This container is running in Hyper-V isolation, which is not supported. No action taken.")
		return
	}

	osName := getContainerOs(ctx, cli, id)
	log.WithFields(log.Fields{
		"os": osName,
	}).Info()

	log.Info("Copying files into container...")
	err := cli.CopyToContainer(ctx, id, "/", osData[osName].Archive, types.CopyToContainerOptions{})
	checkIfError(err)

	log.Info("Running bootstrap script...")
	output, err := execInContainer(ctx, cli, id, osData[osName].Command)
	checkIfError(err)

	fmt.Println("==== Output bootstrap script ====")
	fmt.Println(output)
	fmt.Println("=================================")
	log.Info("Done!")
}

func main() {
	if len(os.Args) != 2 {
		log.Error("No certificate file specified!")
		os.Exit(1)
	}

	cert, err := ioutil.ReadFile(os.Args[1])
	checkIfError(err)

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	checkIfError(err)

	ctx := context.Background()
	checkIfError(err)

	osData, err := getOsData(cert)
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
				bootstrap(ctx, cli, msg.ID, msg.From, osData)
			}
		}
	}
}
