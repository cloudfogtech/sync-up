package docker

import (
	"context"
	"errors"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

func NewDockerCli() (*Docker, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Errorf("Create Docker Cli error: %v", err)
		return nil, err
	}
	return &Docker{
		cli: cli,
	}, nil
}

func (d *Docker) Close() error {
	err := d.cli.Close()
	if err != nil {
		log.Errorf("Close Docker Cli error: %v", err)
		return err
	}
	return nil
}

func (d *Docker) GetContainerList() ([]types.Container, error) {
	containers, err := d.cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Errorf("Docker Cli ContainerList error: %v", err)
		return nil, err
	}
	for i, container := range containers {
		for j, name := range container.Names {
			containers[i].Names[j] = strings.Trim(name, "/")
		}
	}
	log.Debugf("Containers Count: %d", len(containers))
	return containers, err
}

func (d *Docker) GetContainerByName(name string) (types.Container, error) {
	containers, err := d.GetContainerList()
	if err != nil {
		return types.Container{}, err
	}
	for _, container := range containers {
		for _, containerName := range container.Names {
			if containerName == name {
				return container, nil
			}
		}
	}
	return types.Container{}, errors.New("No Container with name=" + name)
}
