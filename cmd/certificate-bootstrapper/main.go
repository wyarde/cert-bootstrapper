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
	"runtime/debug"

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

// func Generate(input ...string) (io.Reader, error) {
// 	files := parseStringPairs(input...)
// 	buf := new(bytes.Buffer)
// 	tw := tar.NewWriter(buf)
// 	for _, file := range files {
// 		name, content := file[0], file[1]
// 		hdr := &tar.Header{
// 			Name: name,
// 			Size: int64(len(content)),
// 		}
// 		if err := tw.WriteHeader(hdr); err != nil {
// 			return nil, err
// 		}
// 		if _, err := tw.Write([]byte(content)); err != nil {
// 			return nil, err
// 		}
// 	}
// 	if err := tw.Close(); err != nil {
// 		return nil, err
// 	}
// 	return buf, nil

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

//go:embed bin/agent-Linux-x86_64
var agentLinux []byte

//go:embed bin/agent-Windows-x86_64.exe
var agentWindows []byte

// OsData type describes Osspecific data requiredfor bootstrapping
type OsData struct {
	Archive io.Reader
	Command []string
}

func getOsData(containerOs string, cert []byte) (osData OsData, err error) {
	certFile := File{
		Name:    "cert.pem",
		Content: cert,
		Mode:    444,
	}

	var (
		agentFile File
		command   []string
	)

	switch containerOs {
	case "linux":
		agentFile = File{
			Name:    "bootstrap-agent",
			Content: agentLinux,
			Mode:    555,
		}
		command = []string{"./bootstrap-agent"}
	case "windows":
		agentFile = File{
			Name:    "bootstrap-agent.exe",
			Content: agentWindows,
			Mode:    555,
		}
		command = []string{"./bootstrap-agent"}

	default:
		log.Error("Unknown operating system: ", containerOs)
		os.Exit(1)
	}

	osData = OsData{
		Archive: generateArchive(agentFile, certFile),
		Command: command,
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
		"id":   id[0:12],
		"from": from,
	}).Info("Bootstrapping container...")

	if isHyperVContainer(ctx, cli, id) {
		log.Warning("This container is running in Hyper-V isolation, which is not supported. No action taken.")
		return
	}

	containerOs := getContainerOs(ctx, cli, id)
	log.WithFields(log.Fields{
		"os": containerOs,
	}).Info()

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

	log.Info("Listening for new containers...")

	for {
		select {
		case err := <-errs:
			log.Error(err)
			os.Exit(1)
		case msg := <-msgs:
			if msg.Status == "start" {
				bootstrap(ctx, cli, msg.ID, msg.From, cert)
			}
		}
	}
}
