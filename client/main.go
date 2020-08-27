package main

import (
	"context"
	"io"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/aclisp/sims/proto"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
)

const (
	// TargetServiceName is the target SIMS server
	TargetServiceName = "go.micro.srv.sims"
)

func main() {
	userID := os.Getenv("USER_ID")
	if userID == "" {
		userID = "homerhuang"
	}
	header := &proto.Header{
		UserId: userID,
	}
	ctx := context.Background()
	service := micro.NewService()
	hub := proto.NewHubService(TargetServiceName, service.Client())
	streamer := proto.NewStreamerService(TargetServiceName, service.Client())

	_, err := hub.Connect(ctx, &proto.ConnectRequest{
		Header: header,
	})
	if err != nil {
		logger.Fatalf("register error: %v", err)
	}

	go func() {
		for range time.Tick(5 * time.Second) {
			_, err := hub.Heartbeat(ctx, &proto.HeartbeatRequest{
				Header: header,
			})
			if err != nil {
				logger.Fatalf("heartbeat error: %v", err)
			}
		}
	}()

	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-term
		hub.Disconnect(ctx, &proto.DisconnectRequest{
			Header: header,
		})
	}()

	header.RequestId = strconv.FormatInt(time.Now().Unix(), 10)
	stream, err := streamer.Events(ctx, &proto.EventsRequest{
		Header: header,
	})
	if err != nil {
		logger.Fatalf("start event stream error: %v", err)
	}
	logger.Infof("[%v] receiving events", header)
	for {
		event, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Fatalf("recv event stream error: %v", err)
		}
		switch event.Type {
		case proto.EventType_EVT_HEARTBEAT:
		default:
			logger.Infof("got event: %v", event)
		}
	}
	logger.Info("done")
}
