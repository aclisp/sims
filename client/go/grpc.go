package im

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	proto "github.com/aclisp/sims/proto/go"
	"google.golang.org/grpc"
)

// GRPCClient TODO
type GRPCClient struct {
	Target string
	UserID string

	subscribeCtx context.Context
	cancel       context.CancelFunc
	conn         *grpc.ClientConn
}

// Unicast TODO
func (c *GRPCClient) Unicast(toUserID, text string) error {
	node := proto.NewPublisherClient(c.conn)
	_, err := node.Unicast(context.TODO(), &proto.UnicastRequest{
		UserId: toUserID,
		Event: &proto.Event{
			Type: proto.EventType_EVT_TEXT,
			Data: []byte(text),
		},
	})
	if err != nil {
		return fmt.Errorf("sims unicast: %w", err)
	}
	return nil
}

// Subscribe TODO
func (c *GRPCClient) Subscribe(callback func(*proto.Event)) {
	c.subscribeCtx, c.cancel = context.WithCancel(context.Background())
	go func() {
		for {
			if err := c.SubscribeEvent(c.subscribeCtx, callback); err != nil {
				log.Printf("subscribe event failure, retrying: %v", err)
			}
			if c.subscribeCtx.Err() != nil {
				break
			}
		}
		log.Println("subscribe done")
	}()
}

// SubscribeEvent TODO
func (c *GRPCClient) SubscribeEvent(ctx context.Context, callback func(*proto.Event)) error {
	conn, err := grpc.DialContext(ctx, c.Target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("grpc dial: %w", err)
	}
	c.conn = conn
	defer conn.Close()

	hub := proto.NewHubClient(conn)
	streamer := proto.NewStreamerClient(conn)

	header := &proto.Header{
		UserId: c.UserID,
	}
	if _, err := hub.Connect(ctx, &proto.ConnectRequest{
		Header: header,
	}); err != nil {
		return fmt.Errorf("node connect: %w", err)
	}

	errHeartbeat := make(chan error, 1)
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			_, err := hub.Heartbeat(ctx, &proto.HeartbeatRequest{
				Header: header,
			})
			if err != nil {
				ticker.Stop()
				errHeartbeat <- err
				return
			}
		}
	}()

	errEvent := make(chan error, 1)
	go func() {
		header := &proto.Header{
			UserId:    c.UserID,
			RequestId: strconv.FormatInt(time.Now().Unix(), 10),
		}
		stream, err := streamer.Events(context.Background(), &proto.EventsRequest{
			Header: header,
		})
		if err != nil {
			errEvent <- err
			return
		}
		for {
			event, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				errEvent <- err
				return
			}
			switch event.Type {
			case proto.EventType_EVT_HEARTBEAT:
			default:
				callback(event)
			}
		}
		close(errEvent)
	}()

	select {
	case err := <-errHeartbeat:
		return fmt.Errorf("node heartbeat: %w", err)
	case err, ok := <-errEvent:
		if ok {
			return fmt.Errorf("node event stream: %w", err)
		}
	}
	return nil //fmt.Errorf("node EOF: %w", io.EOF)
}

// Close TODO
func (c *GRPCClient) Close() error {
	defer func() {
		if c.cancel != nil {
			c.cancel()
		}
	}()

	if c.conn == nil {
		return nil
	}

	hub := proto.NewHubClient(c.conn)
	header := &proto.Header{
		UserId: c.UserID,
	}
	if _, err := hub.Disconnect(context.TODO(), &proto.DisconnectRequest{
		Header: header,
	}); err != nil {
		return fmt.Errorf("node disconnect: %w", err)
	}
	return nil
}
