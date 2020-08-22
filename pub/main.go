package main

import (
	"context"

	"github.com/aclisp/sims/proto"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
)

func main() {
	service := micro.NewService()
	cl := proto.NewPublisherService("go.micro.srv.sims", service.Client())
	_, err := cl.Publish(context.Background(), &proto.PublishRequest{
		UserId: "homerhuang",
		Event: &proto.Event{
			Type: proto.EventType_EVT_TEXT,
			Data: []byte("hello world!"),
		},
	})
	if err != nil {
		logger.Fatalf("publish error: %v", err)
	}
	logger.Info("done")
}
