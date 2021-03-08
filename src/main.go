package main

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"

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
	c := types.ExecConfig{AttachStdin: true, AttachStdout: true, AttachStderr: true, Tty: true, Privileged: true, Cmd: commands}
	execID, err := cli.ContainerExecCreate(ctx, id, c)

	if err != nil {
		return "", err
	}

	res, err := cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return "", err
	}

	s := bufio.NewScanner(res.Reader)
	for s.Scan() {
		// Skip first 8 bytes which contain the header. See https://docs.docker.com/engine/api/v1.24/#attach-to-a-container
		o := s.Text()[8:]
		output += fmt.Sprintf("  > %s\n", o)
	}

	return output, nil
}

//go:embed scripts/bootstrap.sh
var bootstrapScript string

func bootstrap(ctx context.Context, cli *client.Client, id string, from string, cert []byte) {
	log.WithFields(log.Fields{
		"id":   id[0:12],
		"from": from,
	}).Info("Bootstrapping container...")

	reader, err := archive.Generate("bootstrap.sh", bootstrapScript, "cert.crt", string(cert))
	checkIfError(err)

	err = cli.CopyToContainer(ctx, id, "/", reader, types.CopyToContainerOptions{})
	checkIfError(err)

	output, err := execInContainer(ctx, cli, id, []string{"sh", "/bootstrap.sh", "/cert.crt"})
	checkIfError(err)

	log.Info("Running bootstrap script:")
	fmt.Println(output)
	log.Info("Done!")
}

func main() {
	cert, err := ioutil.ReadFile(os.Args[1])
	checkIfError(err)

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	checkIfError(err)

	ctx := context.Background()
	checkIfError(err)

	msgs, errs := cli.Events(ctx, types.EventsOptions{})

	for {
		select {
		case err := <-errs:
			print(err)
		case msg := <-msgs:
			if msg.Status == "start" {
				bootstrap(ctx, cli, msg.ID, msg.From, cert)
			}
		}
	}
}
