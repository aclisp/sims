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

func main() {
	client := im.Client{
		Target: "127.0.0.1:18080",
		UserID: "homerhuang",
	}
	client.Subscribe(func(e *proto.Event) {
		log.Println(e)
	})
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	<-term
	client.Close()
	time.Sleep(time.Second)
}
