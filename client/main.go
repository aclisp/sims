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
	cl := proto.NewHubService("go.micro.srv.sims", service.Client())

	_, err := cl.Connect(ctx, &proto.ConnectRequest{
		Header: header,
	})
	if err != nil {
		logger.Fatalf("register error: %v", err)
	}

	go func() {
		for range time.Tick(5 * time.Second) {
			_, err := cl.Heartbeat(ctx, &proto.HeartbeatRequest{
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
		cl.Disconnect(ctx, &proto.DisconnectRequest{
			Header: header,
		})
	}()

	header.RequestId = strconv.FormatInt(time.Now().Unix(), 10)
	stream, err := cl.Events(ctx, &proto.EventsRequest{
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
