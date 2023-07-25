package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/cloudfogtech/sync-up/internal/api"
	"github.com/cloudfogtech/sync-up/internal/common"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func main() {
	var args common.Command
	p := arg.MustParse(&args)
	s, err := api.NewServer()
	if err != nil {
		p.Fail(fmt.Sprintf("server start fail: %s", err.Error()))
	}
	httpServer, err := s.Serve()
	if err != nil {
		p.Fail(fmt.Sprintf("server start fail: %s", err.Error()))
	}
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})
	if err := g.Wait(); err != nil {
		log.Errorf("server start failed: %v", err)
		p.Fail("server start failed")
	}
}
