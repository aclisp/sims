package main

import (
	"github.com/aclisp/sims/proto"
	"github.com/micro/go-micro/v2/errors"
)

func errorNotRegistered(uid UniqueID) error {
	return errors.BadRequest(proto.ErrorCode_ERR_NOT_FOUND.String(), "not registered for %v", uid)
}

func errorNoConsumer(uid UniqueID) error {
	return errors.InternalServerError(proto.ErrorCode_ERR_NO_CONSUMER.String(), "no consumer for %v", uid)
}
