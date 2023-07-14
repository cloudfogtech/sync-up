package docker

import "github.com/docker/docker/client"

type Docker struct {
	cli *client.Client
}

type ResultChan struct {
	Data chan rune
	Err  chan error
	Exit chan bool
}

func NewResultChan() *ResultChan {
	return &ResultChan{
		Data: make(chan rune),
		Err:  make(chan error),
		Exit: make(chan bool),
	}
}

func (c *ResultChan) Close() {
	close(c.Data)
	close(c.Err)
	close(c.Exit)
}
