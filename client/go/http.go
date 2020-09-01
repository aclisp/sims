package im

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	proto "github.com/aclisp/sims/proto/go"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/golang/protobuf/jsonpb"
	pb "github.com/golang/protobuf/proto"
)

var (
	jsonMarshaler = jsonpb.Marshaler{
		OrigName: true,
	}
	jsonUnmarshaler = jsonpb.Unmarshaler{}
)

// HTTPClient TODO
type HTTPClient struct {
	Target string
	UserID string

	subscribeCtx context.Context
	cancel       context.CancelFunc
	httpClient   http.Client
	wsDialer     ws.Dialer
}

// EventHandler handles server-sent events
type EventHandler interface {
	OnEvent(*proto.Event)
}

// EventHandlerFunc is an adapter to allow the use of ordinary functions as event handlers
type EventHandlerFunc func(*proto.Event)

// OnEvent calls f(e)
func (f EventHandlerFunc) OnEvent(e *proto.Event) {
	f(e)
}

func jsonMarshal(m pb.Message) ([]byte, error) {
	b := new(bytes.Buffer)
	err := jsonMarshaler.Marshal(b, m)
	return b.Bytes(), err
}

func jsonUnmarshal(data []byte, m pb.Message) error {
	return jsonUnmarshaler.Unmarshal(bytes.NewReader(data), m)
}

func (c *HTTPClient) post(ctx context.Context, url, contentType string, body io.Reader) error {
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.New(string(data))
	}
	return nil
}

// Unicast TODO
func (c *HTTPClient) Unicast(toUserID, text string) error {
	buf, _ := jsonMarshal(&proto.UnicastRequest{
		UserId: toUserID,
		Event: &proto.Event{
			Type: proto.EventType_EVT_TEXT,
			Data: []byte(text),
		},
	})
	unicastURL := fmt.Sprintf("http://%s/sims/publisher/unicast", c.Target)
	contentJSON := "application/json"
	if err := c.post(context.TODO(), unicastURL, contentJSON, bytes.NewReader(buf)); err != nil {
		return fmt.Errorf("sims unicast: %w", err)
	}
	return nil
}

// Subscribe TODO
func (c *HTTPClient) Subscribe(h EventHandler) {
	c.subscribeCtx, c.cancel = context.WithCancel(context.Background())
	go func() {
		for {
			if err := c.SubscribeEvent(c.subscribeCtx, h); err != nil {
				log.Printf("subscribe event failure, retrying: %v", err)
			}
			if c.subscribeCtx.Err() != nil {
				break
			}
			time.Sleep(5 * time.Second)
		}
		log.Println("subscribe done")
	}()
}

// SubscribeEvent TODO
func (c *HTTPClient) SubscribeEvent(ctx context.Context, h EventHandler) error {
	var (
		contentJSON  = "application/json"
		connectURL   = fmt.Sprintf("http://%s/sims/hub/connect", c.Target)
		heartbeatURL = fmt.Sprintf("http://%s/sims/hub/heartbeat", c.Target)
		eventsURL    = fmt.Sprintf("ws://%s/sims/streamer/events", c.Target)
		header       = &proto.Header{
			UserId: c.UserID,
		}
		connectReq   = &proto.ConnectRequest{Header: header}
		heartbeatReq = &proto.HeartbeatRequest{Header: header}
		cancel       context.CancelFunc
	)

	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	buf, _ := jsonMarshal(connectReq)
	if err := c.post(ctx, connectURL, contentJSON, bytes.NewReader(buf)); err != nil {
		return fmt.Errorf("node connect: %w", err)
	}

	errHeartbeat := make(chan error, 1)
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			buf, _ := jsonMarshal(heartbeatReq)
			if err := c.post(ctx, heartbeatURL, contentJSON, bytes.NewReader(buf)); err != nil {
				ticker.Stop()
				errHeartbeat <- err
				return
			}
		}
	}()

	errEvent := make(chan error, 1)
	go func() {
		conn, _, _, err := c.wsDialer.Dial(context.Background(), eventsURL)
		if err != nil {
			errEvent <- err
			return
		}
		defer conn.Close()

		buf, _ := jsonMarshal(&proto.EventsRequest{
			Header: &proto.Header{
				UserId:    c.UserID,
				RequestId: strconv.FormatInt(time.Now().Unix(), 10),
			},
		})
		if err := wsutil.WriteClientText(conn, buf); err != nil {
			errEvent <- err
			return
		}

		for {
			data, err := wsutil.ReadServerText(conn)
			if err == io.EOF {
				break
			}
			if err != nil {
				errEvent <- err
				return
			}

			event := new(proto.Event)
			if err := jsonUnmarshal(data, event); err != nil {
				errEvent <- err
				return
			}

			switch event.Type {
			case proto.EventType_EVT_HEARTBEAT:
			default:
				h.OnEvent(event)
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
func (c *HTTPClient) Close() error {
	defer func() {
		if c.cancel != nil {
			c.cancel()
		}
	}()

	buf, _ := jsonMarshal(&proto.DisconnectRequest{
		Header: &proto.Header{
			UserId: c.UserID,
		},
	})
	disconnectURL := fmt.Sprintf("http://%s/sims/hub/disconnect", c.Target)
	contentJSON := "application/json"
	if err := c.post(context.TODO(), disconnectURL, contentJSON, bytes.NewReader(buf)); err != nil {
		return fmt.Errorf("node disconnect: %w", err)
	}
	return nil
}
