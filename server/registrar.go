package main

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/aclisp/sims/proto"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
)

// Registrar TODO
type Registrar struct {
	lock     sync.Mutex
	channels map[UniqueID]*Channel
}

// NewRegistrar TODO
func NewRegistrar() *Registrar {
	return &Registrar{
		channels: make(map[UniqueID]*Channel),
	}
}

// ListChannels TODO
func (reg *Registrar) ListChannels() []*proto.Channel {
	type ch struct {
		UniqueID
		*Channel
	}

	reg.lock.Lock()
	ca := make([]ch, 0, len(reg.channels))
	for uid, channel := range reg.channels {
		ca = append(ca, ch{
			UniqueID: uid,
			Channel:  channel,
		})
	}
	reg.lock.Unlock()

	sort.Slice(ca, func(i, j int) bool { return ca[i].LastHeartbeat.After(ca[j].LastHeartbeat) })

	cb := make([]*proto.Channel, len(ca))
	for i := range ca {
		cb[i] = &proto.Channel{
			UserId:        ca[i].UserID,
			Birth:         ca[i].Birth.Format(time.RFC3339),
			LastHeartbeat: ca[i].LastHeartbeat.Format(time.RFC3339),
			Active:        int32(ca[i].Active.Load()),
		}
	}
	return cb
}

func (reg *Registrar) housekeep() {
	reg.lock.Lock()
	defer reg.lock.Unlock()

	deadline := time.Now().Add(-ChannelInactivity)
	for uid, channel := range reg.channels {
		if channel.LastHeartbeat.Before(deadline) {
			close(channel.EventQueue)
			delete(reg.channels, uid)
		}
	}
}

func (reg *Registrar) close() {
	reg.lock.Lock()
	defer reg.lock.Unlock()

	for _, channel := range reg.channels {
		close(channel.EventQueue)
	}
}

func (reg *Registrar) findEventQueue(uid UniqueID) chan *proto.Event {
	reg.lock.Lock()
	defer reg.lock.Unlock()
	if channel, ok := reg.channels[uid]; ok {
		return channel.EventQueue
	}
	return nil
}

func (reg *Registrar) findChannel(uid UniqueID) *Channel {
	reg.lock.Lock()
	defer reg.lock.Unlock()
	return reg.channels[uid]
}

func (reg *Registrar) heartbeat(uid UniqueID) {
	reg.lock.Lock()
	defer reg.lock.Unlock()
	if channel, ok := reg.channels[uid]; ok {
		channel.LastHeartbeat = time.Now()
	}
}

func (reg *Registrar) createEventQueue(uid UniqueID) {
	reg.lock.Lock()
	defer reg.lock.Unlock()
	if _, ok := reg.channels[uid]; ok {
		return
	}
	channel := &Channel{
		EventQueue:    make(chan *proto.Event),
		Birth:         time.Now(),
		LastHeartbeat: time.Now(),
	}
	reg.channels[uid] = channel
}

func (reg *Registrar) deleteEventQueue(uid UniqueID) {
	reg.lock.Lock()
	defer reg.lock.Unlock()
	if channel, ok := reg.channels[uid]; ok {
		close(channel.EventQueue)
		delete(reg.channels, uid)
	}
}

// Heartbeat TODO
func (reg *Registrar) Heartbeat(ctx context.Context, req *proto.HeartbeatRequest, res *proto.HeartbeatResponse) error {
	uid, err := uniqueIDFromHeader(req.Header)
	if err != nil {
		return err
	}

	events := reg.findEventQueue(uid)
	if events == nil {
		return errorNotRegistered(uid)
	}

	reg.heartbeat(uid)

	event := &proto.Event{Type: proto.EventType_EVT_HEARTBEAT}
	select {
	case events <- event:
	default:
		return errorNoConsumer(uid)
	}
	return nil
}

// Events TODO
func (reg *Registrar) Events(ctx context.Context, req *proto.EventsRequest, stream proto.Streamer_EventsStream) error {
	trace := req.GetHeader().GetRequestId()
	// get notice message queue by user
	uid, err := uniqueIDFromHeader(req.Header)
	if err != nil {
		return err
	}

	channel := reg.findChannel(uid)
	if channel == nil {
		return errorNotRegistered(uid)
	}
	channel.Active.Inc()
	defer channel.Active.Dec()

	// handle event
	logger.Debugf("[%v %v] handling events", uid, trace)
	for event := range channel.EventQueue {
		if err := stream.Send(event); err != nil {
			logger.Errorf("[%v %v] send event to stream error: %v", uid, trace, err)
			return err
		}
	}
	logger.Debugf("[%v %v] no more events", uid, trace)
	return nil
}

// Connect TODO
func (reg *Registrar) Connect(ctx context.Context, req *proto.ConnectRequest, res *proto.ConnectResponse) error {
	// local
	uid, err := uniqueIDFromHeader(req.Header)
	if err != nil {
		return err
	}
	reg.createEventQueue(uid)
	// persist: which server box the uid belongs to?
	serverName := gService.Server().Options().Name
	serverID := gService.Server().Options().Id
	myNodeID := serverName + "-" + serverID
	services, err := gService.Options().Registry.GetService(serverName)
	if err != nil {
		logger.Errorf("get service %q from registry: %v", serverName, err)
		return err
	}
	var myNode *registry.Node
	for _, service := range services {
		for _, node := range service.Nodes {
			if myNodeID == node.Id {
				myNode = node
				break
			}
		}
	}
	if myNode == nil {
		err := errors.New("self node not found in registry")
		logger.Errorf("get service %q from registry: %v", serverName, err)
		return err
	}
	logger.Debugf("%v connected to %v", uid, myNode.Address)
	return nil
}

// Disconnect TODO
func (reg *Registrar) Disconnect(ctx context.Context, req *proto.DisconnectRequest, res *proto.DisconnectResponse) error {
	uid, err := uniqueIDFromHeader(req.Header)
	if err != nil {
		return err
	}
	reg.deleteEventQueue(uid)
	return nil
}

// List TODO
func (reg *Registrar) List(ctx context.Context, req *proto.ListRequest, res *proto.ListResponse) error {
	res.Channels = reg.ListChannels()
	return nil
}
