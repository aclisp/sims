package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	im "github.com/aclisp/sims/client/go"
	proto "github.com/aclisp/sims/proto/go"
)

var (
	connectionInstanceCount = flag.Int("c", 1, "the count of IM connection instance")
	userIdentityPrefix      = flag.String("u", "user", "the prefix of user identity. The first user would be named as `user_1`")
	targetAPIServerAddress  = flag.String("t", "127.0.0.1:8080", "the API server address to connect to")
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
	flag.Parse()
	conns := make([]*Connection, *connectionInstanceCount)
	for i := range conns {
		conn := &Connection{
			HTTPClient: im.HTTPClient{
				Target: *targetAPIServerAddress,
				UserID: fmt.Sprintf("%s_%d", *userIdentityPrefix, i+1),
			},
		}
		conn.Subscribe(conn)
		conns[i] = conn
	}
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	<-term
	for _, conn := range conns {
		conn.Close()
	}
	time.Sleep(time.Second)
}
