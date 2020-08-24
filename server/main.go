package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"sort"
	"sync"
	"time"

	"github.com/aclisp/sims/proto"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/logger"
	"google.golang.org/grpc"
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
func (reg *Registrar) Events(ctx context.Context, req *proto.EventsRequest, stream proto.Hub_EventsStream) error {
	trace := req.GetHeader().GetRequestId()
	// get notice message queue by user
	uid, err := uniqueIDFromHeader(req.Header)
	if err != nil {
		return err
	}
	events := reg.findEventQueue(uid)
	if events == nil {
		return errorNotRegistered(uid)
	}
	// handle event
	logger.Infof("[%v %v] handling events", uid, trace)
	for event := range events {
		if err := stream.Send(event); err != nil {
			logger.Errorf("[%v %v] send event to stream error: %v", uid, trace, err)
			return err
		}
	}
	logger.Infof("[%v %v] no more events", uid, trace)
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
	// persist
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
	if req.Event == nil {
		return errors.BadRequest(proto.ErrorCode_ERR_MISSING_EVENT.String(), "nil event for %v", uid)
	}
	if req.Event.Type == proto.EventType_EVT_HEARTBEAT {
		return errors.BadRequest(proto.ErrorCode_ERR_INVALID_EVENT_TYPE.String(), "event type should not be EVT_HEARTBEAT")
	}
	event := req.Event
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
		micro.BeforeStop(func() error {
			gRegistrar.close()
			return nil
		}),
		micro.Flags(&cli.StringFlag{
			Name:    "pprof_address",
			EnvVars: []string{"PPROF_ADDRESS"},
			Usage:   "Bind address for pprof and grpc.EnableTracing. 127.0.0.1:6060",
		}),
		micro.Action(func(ctx *cli.Context) error {
			if addr := ctx.String("pprof_address"); len(addr) > 0 {
				// for pprof and trace
				grpc.EnableTracing = true
				go func() { logger.Warn(http.ListenAndServe(addr, nil)) }()
			}
			return nil
		}),
	)

	service.Init()

	proto.RegisterHubHandler(service.Server(), gRegistrar)
	proto.RegisterPublisherHandler(service.Server(), new(Publisher))

	logger.Info("run")
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			gRegistrar.housekeep()
		}
	}()
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
	ticker.Stop()
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

func errorNotRegistered(uid UniqueID) error {
	return errors.BadRequest(proto.ErrorCode_ERR_NOT_FOUND.String(), "not registered for %v", uid)
}

func errorNoConsumer(uid UniqueID) error {
	return errors.InternalServerError(proto.ErrorCode_ERR_NO_CONSUMER.String(), "no consumer for %v", uid)
}
