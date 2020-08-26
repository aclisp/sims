package main

import (
	"github.com/aclisp/sims/proto"
	"github.com/micro/go-micro/v2/errors"
)

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
