package main

import (
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/aclisp/sims/proto"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	"google.golang.org/grpc"
)

const (
	// MicroServiceName is the name of the service
	MicroServiceName = "go.micro.srv.sims"
	// HousekeepInterval is the duration after which we do housekeeping
	HousekeepInterval = 5 * time.Second
	// ChannelInactivity is the duration after which an inactive channel is closed by the server
	ChannelInactivity = 10 * time.Second
)

var (
	gRegistrar *Registrar = NewRegistrar()
)

func main() {
	service := micro.NewService(
		micro.Name(MicroServiceName),
		micro.BeforeStop(func() error {
			gRegistrar.close()
			return nil
		}),
		micro.Flags(&cli.StringFlag{
			Name:    "pprof_address",
			EnvVars: []string{"PPROF_ADDRESS"},
			Usage:   "Bind address for pprof and grpc.EnableTracing. 127.0.0.1:6060",
		}),
		micro.Action(func(ctx *cli.Context) error {
			if addr := ctx.String("pprof_address"); len(addr) > 0 {
				// for pprof and trace
				grpc.EnableTracing = true
				go func() { logger.Warn(http.ListenAndServe(addr, nil)) }()
			}
			return nil
		}),
	)

	service.Init()

	proto.RegisterHubHandler(service.Server(), gRegistrar)
	proto.RegisterStreamerHandler(service.Server(), gRegistrar)
	proto.RegisterPublisherHandler(service.Server(), new(Publisher))

	logger.Info("run")
	ticker := time.NewTicker(HousekeepInterval)
	go func() {
		for range ticker.C {
			gRegistrar.housekeep()
		}
	}()
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
	ticker.Stop()
}
