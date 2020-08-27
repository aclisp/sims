package main

import (
	"context"

	"github.com/aclisp/sims/proto"
	"github.com/micro/go-micro/v2/errors"
)

// Publisher TODO
type Publisher struct{}

// Unicast TODO
func (pub *Publisher) Unicast(ctx context.Context, req *proto.UnicastRequest, res *proto.UnicastResponse) error {
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

// Multicast TODO
func (pub *Publisher) Multicast(ctx context.Context, req *proto.MulticastRequest, res *proto.MulticastResponse) error {
	if len(req.UserId) == 0 {
		return errors.BadRequest(proto.ErrorCode_ERR_MISSING_USERID.String(), "need at least one user_id")
	}
	res.UserErrcode = make(map[string]proto.ErrorCode)
	for _, u := range req.UserId {
		if err := pub.Unicast(ctx, &proto.UnicastRequest{
			UserId:       u,
			Event:        req.Event,
			UserSelector: req.UserSelector[u],
		}, &proto.UnicastResponse{}); err != nil {
			var errCode proto.ErrorCode
			if ierr, ok := err.(*errors.Error); ok {
				errCode = proto.ErrorCode(proto.ErrorCode_value[ierr.Id])
			} else {
				errCode = proto.ErrorCode_ERR_UNSPECIFIED
			}
			res.UserErrcode[u] = errCode
		}
	}
	return nil
}
