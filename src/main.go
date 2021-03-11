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

func bootstrap(ctx context.Context, cli *client.Client, id string, from string, cert []byte) {
	log.WithFields(log.Fields{
		"id":   id[0:12],
		"from": from,
	}).Info("Bootstrapping container...")

	info, err := cli.ContainerInspect(ctx, id)

	if info.ContainerJSONBase.HostConfig.Isolation == "hyperv" {
		log.Warning("This container is running in Hyper-V isolation, which is not supported. No action taken.")
		return
	}

	log.WithFields(log.Fields{
		"os": info.ContainerJSONBase.Platform,
	}).Info()

	var reader io.Reader
	var bootstrapCmd []string

	switch info.ContainerJSONBase.Platform {
	case "linux":
		reader, err = archive.Generate("bootstrap.sh", bootstrapLinux, "cert.crt", string(cert))
		checkIfError(err)
		bootstrapCmd = []string{"sh", "/bootstrap.sh", "/cert.crt"}
	case "windows":
		reader, err = archive.Generate("bootstrap.ps1", bootstrapWindows, "cert.crt", string(cert))
		checkIfError(err)
		bootstrapCmd = []string{"powershell", "-Command", "/bootstrap.ps1", "/cert.crt"}
	default:
		log.Warning("Don't know about this operating system. No action taken.")
		return
	}

	log.Info("Copying files into container...")
	err = cli.CopyToContainer(ctx, id, "/", reader, types.CopyToContainerOptions{})
	checkIfError(err)

	log.Info("Running bootstrap script...")
	output, err := execInContainer(ctx, cli, id, bootstrapCmd)
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

	msgs, errs := cli.Events(ctx, types.EventsOptions{})
	checkIfError(err)

	log.Info("Listening for new containers...")

	for {
		select {
		case err := <-errs:
			print(err)
			os.Exit(1)
		case msg := <-msgs:
			if msg.Status == "start" {
				bootstrap(ctx, cli, msg.ID, msg.From, cert)
			}
		}
	}
}
