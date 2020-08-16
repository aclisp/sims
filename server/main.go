package main

import (
	"context"
	"encoding/json"
	"sort"
	"sync"
	"time"

	"github.com/aclisp/sims/proto"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/metadata"
)

var (
	gRegistrar *Registrar = NewRegistrar()
)

// Channel TODO
type Channel struct {
	EventQueue    chan *proto.Event
	Birth         time.Time
	LastHeartbeat time.Time
}

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
	reg.lock.Lock()
	defer reg.lock.Unlock()

	type ch struct {
		UniqueID
		*Channel
	}
	ca := make([]ch, 0, len(reg.channels))
	for uid, channel := range reg.channels {
		ca = append(ca, ch{
			UniqueID: uid,
			Channel:  channel,
		})
	}
	sort.Slice(ca, func(i, j int) bool { return ca[i].LastHeartbeat.After(ca[j].LastHeartbeat) })

	cb := make([]*proto.Channel, len(ca))
	for i := range ca {
		cb[i] = &proto.Channel{
			UserId:        ca[i].UserID,
			Birth:         ca[i].Birth.Format(time.RFC3339),
			LastHeartbeat: ca[i].LastHeartbeat.Format(time.RFC3339),
		}
	}
	return cb
}

func (reg *Registrar) housekeep() {
	reg.lock.Lock()
	defer reg.lock.Unlock()

	deadline := time.Now().Add(-10 * time.Second)
	for uid, channel := range reg.channels {
		if channel.LastHeartbeat.Before(deadline) {
			close(channel.EventQueue)
			delete(reg.channels, uid)
		}
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

// Heartbeat TODO
func (reg *Registrar) Heartbeat(ctx context.Context, req *proto.HeartbeatRequest, res *proto.HeartbeatResponse) error {
	md, _ := metadata.FromContext(ctx)
	header := metadataToHeader(md)
	uid, err := uniqueIDFromHeader(header)
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

// EventStream TODO
func (reg *Registrar) EventStream(ctx context.Context, req *proto.EventStreamRequest, stream proto.Registrar_EventStreamStream) error {
	md, _ := metadata.FromContext(ctx)
	header := metadataToHeader(md)
	trace := header.RequestId
	// get notice message queue by user
	uid, err := uniqueIDFromHeader(header)
	if err != nil {
		return err
	}
	events := reg.findEventQueue(uid)
	if events == nil {
		return errorNotRegistered(uid)
	}
	errch := make(chan error)
	// handle event
	go func() {
		logger.Infof("[%v %v] handling events", uid, trace)
		for event := range events {
			if err := stream.Send(event); err != nil {
				errch <- err
				return
			}
		}
		close(errch)
	}()
	if err, ok := <-errch; ok {
		logger.Errorf("[%v %v] send event to stream error: %v", uid, trace, err)
		return err
	}
	logger.Infof("[%v %v] no more events", uid, trace)
	return nil
}

// Register TODO
func (reg *Registrar) Register(ctx context.Context, req *proto.RegisterRequest, res *proto.RegisterResponse) error {
	md, _ := metadata.FromContext(ctx)
	header := metadataToHeader(md)
	// local
	uid, err := uniqueIDFromHeader(header)
	if err != nil {
		return err
	}
	reg.createEventQueue(uid)
	// persist
	return nil
}

// Unregister TODO
func (reg *Registrar) Unregister(context.Context, *proto.UnregisterRequest, *proto.UnregisterResponse) error {
	return nil
}

// List TODO
func (reg *Registrar) List(ctx context.Context, req *proto.ListRequest, res *proto.ListResponse) error {
	res.Channels = reg.ListChannels()
	return nil
}

// Publisher TODO
type Publisher struct{}

// Publish TODO
func (pub *Publisher) Publish(ctx context.Context, req *proto.PublishRequest, res *proto.PublishResponse) error {
	uid := UniqueID{
		UserID: req.UserId,
	}
	events := gRegistrar.findEventQueue(uid)
	if events == nil {
		return errorNotRegistered(uid)
	}
	event := &proto.Event{
		Type: proto.EventType_EVT_TEXT,
		Data: []byte(time.Now().String()),
	}
	select {
	case events <- event:
	default:
		return errorNoConsumer(uid)
	}
	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("go.micro.srv.sims"),
	)

	proto.RegisterRegistrarHandler(service.Server(), gRegistrar)
	proto.RegisterPublisherHandler(service.Server(), new(Publisher))

	logger.Info("run")
	go func() {
		for range time.Tick(5 * time.Second) {
			gRegistrar.housekeep()
		}
	}()
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}

// UniqueID TODO
type UniqueID struct {
	UserID string
}

func uniqueIDFromHeader(header *proto.Header) (UniqueID, error) {
	if header == nil {
		return UniqueID{}, errors.BadRequest(proto.ErrorCode_ERR_MISSING_HEADER.String(), "")
	}
	userID := header.GetUserId()
	if userID == "" {
		return UniqueID{}, errors.BadRequest(proto.ErrorCode_ERR_MISSING_USERID.String(), "")
	}
	return UniqueID{
		UserID: userID,
	}, nil
}

func metadataToHeader(md metadata.Metadata) *proto.Header {
	header := &proto.Header{}
	buf, _ := json.Marshal(md)
	json.Unmarshal(buf, header)
	return header
}

func errorNotRegistered(uid UniqueID) error {
	return errors.BadRequest(proto.ErrorCode_ERR_NOT_FOUND.String(), "not registered for %v", uid)
}

func errorNoConsumer(uid UniqueID) error {
	return errors.InternalServerError(proto.ErrorCode_ERR_NO_CONSUMER.String(), "no consumer for %v", uid)
}
