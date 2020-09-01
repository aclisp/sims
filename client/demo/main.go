package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	im "github.com/aclisp/sims/client/go"
	proto "github.com/aclisp/sims/proto/go"
)

// Connection is a instance of IM connection
type Connection struct {
	im.HTTPClient
}

// OnEvent handles server-sent events
func (c *Connection) OnEvent(e *proto.Event) {
	log.Println(c.UserID, e)
}

func main() {
	conn := &Connection{
		HTTPClient: im.HTTPClient{
			Target: "127.0.0.1:8080",
			UserID: "homerhuang",
		},
	}
	conn.Subscribe(conn)
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	<-term
	conn.Close()
	time.Sleep(time.Second)
}
