package main

import (
	"context"
	"os"

	"github.com/aclisp/sims/proto"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
)

func main() {
	service := micro.NewService()
	cl := proto.NewPublisherService("go.micro.srv.sims", service.Client())
	res, err := cl.Multicast(context.Background(), &proto.MulticastRequest{
		UserId: os.Args[1:],
		Event: &proto.Event{
			Type: proto.EventType_EVT_TEXT,
			Data: []byte("hello world!"),
		},
	})
	if err != nil {
		logger.Fatalf("publish error: %v", err)
	}
	logger.Infof("done: %v", res)
}
