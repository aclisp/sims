package im_test

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	im "github.com/aclisp/sims/client/go"
	proto "github.com/aclisp/sims/proto/go"
)

func TestEventGRPC(t *testing.T) {
	bin := bin()

	server := Command{Path: bin, Name: "server", Args: []string{"--server_address", "127.0.0.1:18080"}}
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}

	client := im.GRPCClient{
		Target: "127.0.0.1:18080",
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

func TestCloseGRPC(t *testing.T) {
	bin := bin()

	server := Command{Path: bin, Name: "server", Args: []string{"--server_address", "127.0.0.1:18080"}}
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}

	client := im.GRPCClient{
		Target: "127.0.0.1:18080",
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

// bin returns the project `bin` dir path; must be called from TestXXX
func bin() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Join(filepath.Dir(filename), "../../bin")
}

type Command struct {
	Path string
	Name string
	Args []string
	cmd  *exec.Cmd
	out  bytes.Buffer
}

func (s *Command) Start() error {
	s.cmd = exec.Command(filepath.Join(s.Path, s.Name), s.Args...)
	s.cmd.Stdout = &s.out
	s.cmd.Stderr = &s.out
	if err := s.cmd.Start(); err != nil {
		return err
	}
	time.Sleep(time.Second)
	return nil
}

func (s *Command) Stop() {
	s.cmd.Process.Signal(os.Interrupt)
	s.cmd.Wait()
}

func (s *Command) Out() []string {
	return strings.Split(s.out.String(), "\n")
}
