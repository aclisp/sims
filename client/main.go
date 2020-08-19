package main

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/aclisp/sims/proto"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/metadata"
)

func headerToMetadata(header *proto.Header) metadata.Metadata {
	var md metadata.Metadata
	buf, _ := json.Marshal(header)
	json.Unmarshal(buf, &md)
	return md
}

func main() {
	userID := os.Getenv("USER_ID")
	if userID == "" {
		userID = "homerhuang"
	}
	header := &proto.Header{
		UserId: userID,
	}
	ctx := metadata.NewContext(context.Background(), headerToMetadata(header))
	service := micro.NewService()
	cl := proto.NewIMNodeService("go.micro.srv.sims", service.Client())

	_, err := cl.Register(ctx, &proto.RegisterRequest{})
	if err != nil {
		logger.Fatalf("register error: %v", err)
	}

	go func() {
		for range time.Tick(5 * time.Second) {
			_, err := cl.Heartbeat(ctx, &proto.HeartbeatRequest{})
			if err != nil {
				logger.Fatalf("heartbeat error: %v", err)
			}
		}
	}()

	header.RequestId = strconv.FormatInt(time.Now().Unix(), 10)
	ctx = metadata.NewContext(context.Background(), headerToMetadata(header))
	stream, err := cl.EventStream(ctx, &proto.EventStreamRequest{})
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
