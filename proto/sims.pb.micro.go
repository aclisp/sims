// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: sims.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Hub service

func NewHubEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Hub service

type HubService interface {
	Connect(ctx context.Context, in *ConnectRequest, opts ...client.CallOption) (*ConnectResponse, error)
	Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...client.CallOption) (*HeartbeatResponse, error)
	Disconnect(ctx context.Context, in *DisconnectRequest, opts ...client.CallOption) (*DisconnectResponse, error)
	List(ctx context.Context, in *ListRequest, opts ...client.CallOption) (*ListResponse, error)
}

type hubService struct {
	c    client.Client
	name string
}

func NewHubService(name string, c client.Client) HubService {
	return &hubService{
		c:    c,
		name: name,
	}
}

func (c *hubService) Connect(ctx context.Context, in *ConnectRequest, opts ...client.CallOption) (*ConnectResponse, error) {
	req := c.c.NewRequest(c.name, "Hub.Connect", in)
	out := new(ConnectResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hubService) Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...client.CallOption) (*HeartbeatResponse, error) {
	req := c.c.NewRequest(c.name, "Hub.Heartbeat", in)
	out := new(HeartbeatResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hubService) Disconnect(ctx context.Context, in *DisconnectRequest, opts ...client.CallOption) (*DisconnectResponse, error) {
	req := c.c.NewRequest(c.name, "Hub.Disconnect", in)
	out := new(DisconnectResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hubService) List(ctx context.Context, in *ListRequest, opts ...client.CallOption) (*ListResponse, error) {
	req := c.c.NewRequest(c.name, "Hub.List", in)
	out := new(ListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Hub service

type HubHandler interface {
	Connect(context.Context, *ConnectRequest, *ConnectResponse) error
	Heartbeat(context.Context, *HeartbeatRequest, *HeartbeatResponse) error
	Disconnect(context.Context, *DisconnectRequest, *DisconnectResponse) error
	List(context.Context, *ListRequest, *ListResponse) error
}

func RegisterHubHandler(s server.Server, hdlr HubHandler, opts ...server.HandlerOption) error {
	type hub interface {
		Connect(ctx context.Context, in *ConnectRequest, out *ConnectResponse) error
		Heartbeat(ctx context.Context, in *HeartbeatRequest, out *HeartbeatResponse) error
		Disconnect(ctx context.Context, in *DisconnectRequest, out *DisconnectResponse) error
		List(ctx context.Context, in *ListRequest, out *ListResponse) error
	}
	type Hub struct {
		hub
	}
	h := &hubHandler{hdlr}
	return s.Handle(s.NewHandler(&Hub{h}, opts...))
}

type hubHandler struct {
	HubHandler
}

func (h *hubHandler) Connect(ctx context.Context, in *ConnectRequest, out *ConnectResponse) error {
	return h.HubHandler.Connect(ctx, in, out)
}

func (h *hubHandler) Heartbeat(ctx context.Context, in *HeartbeatRequest, out *HeartbeatResponse) error {
	return h.HubHandler.Heartbeat(ctx, in, out)
}

func (h *hubHandler) Disconnect(ctx context.Context, in *DisconnectRequest, out *DisconnectResponse) error {
	return h.HubHandler.Disconnect(ctx, in, out)
}

func (h *hubHandler) List(ctx context.Context, in *ListRequest, out *ListResponse) error {
	return h.HubHandler.List(ctx, in, out)
}

// Api Endpoints for Streamer service

func NewStreamerEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Streamer service

type StreamerService interface {
	Events(ctx context.Context, in *EventsRequest, opts ...client.CallOption) (Streamer_EventsService, error)
}

type streamerService struct {
	c    client.Client
	name string
}

func NewStreamerService(name string, c client.Client) StreamerService {
	return &streamerService{
		c:    c,
		name: name,
	}
}

func (c *streamerService) Events(ctx context.Context, in *EventsRequest, opts ...client.CallOption) (Streamer_EventsService, error) {
	req := c.c.NewRequest(c.name, "Streamer.Events", &EventsRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &streamerServiceEvents{stream}, nil
}

type Streamer_EventsService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*Event, error)
}

type streamerServiceEvents struct {
	stream client.Stream
}

func (x *streamerServiceEvents) Close() error {
	return x.stream.Close()
}

func (x *streamerServiceEvents) Context() context.Context {
	return x.stream.Context()
}

func (x *streamerServiceEvents) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *streamerServiceEvents) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *streamerServiceEvents) Recv() (*Event, error) {
	m := new(Event)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Streamer service

type StreamerHandler interface {
	Events(context.Context, *EventsRequest, Streamer_EventsStream) error
}

func RegisterStreamerHandler(s server.Server, hdlr StreamerHandler, opts ...server.HandlerOption) error {
	type streamer interface {
		Events(ctx context.Context, stream server.Stream) error
	}
	type Streamer struct {
		streamer
	}
	h := &streamerHandler{hdlr}
	return s.Handle(s.NewHandler(&Streamer{h}, opts...))
}

type streamerHandler struct {
	StreamerHandler
}

func (h *streamerHandler) Events(ctx context.Context, stream server.Stream) error {
	m := new(EventsRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.StreamerHandler.Events(ctx, m, &streamerEventsStream{stream})
}

type Streamer_EventsStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Event) error
}

type streamerEventsStream struct {
	stream server.Stream
}

func (x *streamerEventsStream) Close() error {
	return x.stream.Close()
}

func (x *streamerEventsStream) Context() context.Context {
	return x.stream.Context()
}

func (x *streamerEventsStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *streamerEventsStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *streamerEventsStream) Send(m *Event) error {
	return x.stream.Send(m)
}

// Api Endpoints for Publisher service

func NewPublisherEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Publisher service

type PublisherService interface {
	Unicast(ctx context.Context, in *UnicastRequest, opts ...client.CallOption) (*UnicastResponse, error)
	Multicast(ctx context.Context, in *MulticastRequest, opts ...client.CallOption) (*MulticastResponse, error)
}

type publisherService struct {
	c    client.Client
	name string
}

func NewPublisherService(name string, c client.Client) PublisherService {
	return &publisherService{
		c:    c,
		name: name,
	}
}

func (c *publisherService) Unicast(ctx context.Context, in *UnicastRequest, opts ...client.CallOption) (*UnicastResponse, error) {
	req := c.c.NewRequest(c.name, "Publisher.Unicast", in)
	out := new(UnicastResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publisherService) Multicast(ctx context.Context, in *MulticastRequest, opts ...client.CallOption) (*MulticastResponse, error) {
	req := c.c.NewRequest(c.name, "Publisher.Multicast", in)
	out := new(MulticastResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Publisher service

type PublisherHandler interface {
	Unicast(context.Context, *UnicastRequest, *UnicastResponse) error
	Multicast(context.Context, *MulticastRequest, *MulticastResponse) error
}

func RegisterPublisherHandler(s server.Server, hdlr PublisherHandler, opts ...server.HandlerOption) error {
	type publisher interface {
		Unicast(ctx context.Context, in *UnicastRequest, out *UnicastResponse) error
		Multicast(ctx context.Context, in *MulticastRequest, out *MulticastResponse) error
	}
	type Publisher struct {
		publisher
	}
	h := &publisherHandler{hdlr}
	return s.Handle(s.NewHandler(&Publisher{h}, opts...))
}

type publisherHandler struct {
	PublisherHandler
}

func (h *publisherHandler) Unicast(ctx context.Context, in *UnicastRequest, out *UnicastResponse) error {
	return h.PublisherHandler.Unicast(ctx, in, out)
}

func (h *publisherHandler) Multicast(ctx context.Context, in *MulticastRequest, out *MulticastResponse) error {
	return h.PublisherHandler.Multicast(ctx, in, out)
}
