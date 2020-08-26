package main

import (
	"time"

	"github.com/aclisp/sims/proto"
	"go.uber.org/atomic"
)

// Channel TODO
type Channel struct {
	EventQueue    chan *proto.Event
	Birth         time.Time
	LastHeartbeat time.Time
	Active        atomic.Uint32
}
