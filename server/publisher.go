package main

import (
	"context"

	"github.com/aclisp/sims/proto"
	"github.com/micro/go-micro/v2/errors"
)

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
