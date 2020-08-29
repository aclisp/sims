package im_test

import (
	"context"
	"testing"
	"time"

	im "github.com/aclisp/sims/client/go"
	proto "github.com/aclisp/sims/proto/go"
)

func TestEventHTTP(t *testing.T) {
	bin := bin()

	server := Command{Path: bin, Name: "server", Args: []string{"--server_address", "127.0.0.1:18080"}}
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}

	client := im.HTTPClient{
		Target: "127.0.0.1:8080",
		UserID: "homerhuang",
	}

	const Text = "hello"

	errSubscribe := make(chan error, 1)
	go func() {
		if err := client.SubscribeEvent(context.Background(), func(e *proto.Event) {
			t.Log(e)
			if e.Type == proto.EventType_EVT_TEXT && string(e.Data) == Text {
			} else {
				t.Fail()
			}
		}); err != nil {
			errSubscribe <- err
		}
		close(errSubscribe)
	}()
	time.Sleep(time.Second)

	if err := client.Unicast("homerhuang", Text); err != nil {
		t.Log(err)
		t.Fail()
	}

	server.Stop()

	if err, ok := <-errSubscribe; ok {
		t.Log(err)
		t.Fail()
	}

	for _, out := range server.Out() {
		t.Log(out)
	}
}

func TestCloseHTTP(t *testing.T) {
	bin := bin()

	server := Command{Path: bin, Name: "server", Args: []string{"--server_address", "127.0.0.1:18080"}}
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}

	client := im.HTTPClient{
		Target: "127.0.0.1:8080",
		UserID: "homerhuang",
	}

	errSubscribe := make(chan error, 1)
	go func() {
		if err := client.SubscribeEvent(context.Background(), func(e *proto.Event) {}); err != nil {
			errSubscribe <- err
		}
		close(errSubscribe)
	}()
	time.Sleep(time.Second)

	if err := client.Close(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if err, ok := <-errSubscribe; ok {
		t.Log(err)
		t.Fail()
	}

	server.Stop()

	for _, out := range server.Out() {
		t.Log(out)
	}
}
