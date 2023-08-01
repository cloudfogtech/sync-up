package docker

import (
	"bufio"
	"context"
	"errors"
	"io"
	"strings"

	"github.com/cloudfogtech/sync-up/internal/utils"
	"github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"
)

func (d *Docker) RunCommandAsync(container types.Container, commands []string, resultChan *ResultChan) (chan bool, chan error, error) {
	execResp, err := d.cli.ContainerExecCreate(context.Background(), container.ID, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          commands,
	})
	if err != nil {
		log.Debugf("Docker Cli ContainerExecCreate error: '%s', %v", strings.Join(commands, " "), err)
		return nil, nil, err
	}
	log.Debugf("ContainerExecCreate - exec id: %s", execResp.ID)
	containerResponse, err := d.cli.ContainerExecAttach(context.Background(), execResp.ID, types.ExecConfig{})
	if err != nil {
		log.Debugf("Docker Cli ContainerExecAttach error: '%s', %v", strings.Join(commands, " "), err)
		return nil, nil, err
	}
	log.Debugf("ContainerExecAttach - containerId: %s, execId: %s", container.ID, execResp.ID)
	bufReader := bufio.NewReader(containerResponse.Reader)
	finishChan := make(chan bool)
	successChan := make(chan error)
	go func() {
		for {
			r, _, err := bufReader.ReadRune()
			if err != nil {
				if err == io.EOF {
					log.Debugf("container response read end")
				} else {
					log.Errorf("container response read error, %v", err)
					resultChan.Err <- err
					successChan <- err
				}
				log.Debugf("container response read exit")
				resultChan.Exit <- true
				break
			}
			resultChan.Data <- r
		}
		containerResponse.Close()
		finishChan <- true
		successChan <- nil
	}()
	log.Debugf("Container Name='%s', ID='%s', Run='%s'", strings.Join(container.Names, ","), container.ID, strings.Join(commands, " "))
	return finishChan, successChan, nil
}

func (d *Docker) RunCommandSync(container types.Container, commands []string) (string, error) {
	createResponse, err := d.cli.ContainerExecCreate(context.Background(), container.ID, types.ExecConfig{
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          commands,
	})
	if err != nil {
		log.Debugf("Docker Cli ContainerExecCreate error: '%s', %v", strings.Join(commands, " "), err)
		return "", err
	}
	log.Debugf("ContainerExecCreate - exec id: %s", createResponse.ID)
	attachResponse, err := d.cli.ContainerExecAttach(context.Background(), createResponse.ID, types.ExecConfig{})
	if err != nil {
		log.Debugf("Docker Cli ContainerExecAttach error: '%s', %v", strings.Join(commands, " "), err)
		return "", err
	}
	log.Debugf("ContainerExecAttach - containerId: %s, execId: %s", container.ID, createResponse.ID)
	result, err := getDataFromReader(attachResponse.Reader)
	if err != nil {
		return "", err
	}
	attachResponse.Close()
	inspectResponse, err := d.cli.ContainerExecInspect(context.Background(), createResponse.ID)
	if err != nil {
		log.Debugf("Docker Cli ContainerExecInspect error: '%s', %v", strings.Join(commands, " "), err)
		return "", err
	}
	if inspectResponse.ExitCode != 0 {
		log.Debugf("Container Exec Error[%d] Name='%s', ID='%s', Run='%s'", inspectResponse.ExitCode, strings.Join(container.Names, ","),
			container.ID, strings.Join(commands, " "))
		return "", errors.New(result)
	}
	log.Debugf("Container Exec Success[%d] Name='%s', ID='%s', Run='%s'", inspectResponse.ExitCode, strings.Join(container.Names, ","), container.ID, strings.Join(commands, " "))
	return result, nil
}

func getDataFromReader(reader *bufio.Reader) (string, error) {
	result := ""
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				log.Debugf("container response read end")
			} else {
				log.Debugf("container response read error, %v", err)
				return "", err
			}
			log.Debugf("container response read exit")
			break
		}
		result += utils.GetPrintable(string(r))
	}
	return result, nil
}
