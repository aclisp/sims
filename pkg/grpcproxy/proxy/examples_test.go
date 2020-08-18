// Copyright 2017-2018 Valient Gough
// Copyright 2017 Michal Witkowski
// All Rights Reserved.
// See LICENSE for licensing terms.

package proxy_test

import (
	"context"
	"strings"

	"github.com/aclisp/sims/pkg/grpcproxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

var (
	director proxy.StreamDirector
)

func ExampleRegisterService() {
	// A gRPC server with the proxying codec enabled.
	server := grpc.NewServer(grpc.CustomCodec(proxy.Codec()))
	// Register a TestService with 4 of its methods explicitly.
	proxy.RegisterService(server, director,
		"vgough.testproto.TestService",
		"PingEmpty", "Ping", "PingError", "PingList")
}

func ExampleTransparentHandler() {
	grpc.NewServer(
		grpc.CustomCodec(proxy.Codec()),
		grpc.UnknownServiceHandler(proxy.TransparentHandler(director)))
}

// Provide sa simple example of a director that shields internal services and dials a staging or production backend.
// This is a *very naive* implementation that creates a new connection on every request. Consider using pooling.
type ExampleDirector struct {
}

func (d *ExampleDirector) Connect(ctx context.Context, method string) (context.Context, *grpc.ClientConn, error) {
	// Make sure we never forward internal services.
	if strings.HasPrefix(method, "/com.example.internal.") {
		return nil, nil, grpc.Errorf(codes.Unimplemented, "Unknown method")
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		// Decide on which backend to dial
		if val, exists := md[":authority"]; exists && val[0] == "staging.api.example.com" {
			// Make sure we use DialContext so the dialing can be cancelled/time out together with the context.
			conn, err := grpc.DialContext(ctx, "api-service.staging.svc.local", grpc.WithCodec(proxy.Codec()))
			return ctx, conn, err
		} else if val, exists := md[":authority"]; exists && val[0] == "api.example.com" {
			conn, err := grpc.DialContext(ctx, "api-service.prod.svc.local", grpc.WithCodec(proxy.Codec()))
			return ctx, conn, err
		}
	}
	return nil, nil, grpc.Errorf(codes.Unimplemented, "Unknown method")
}

func (d *ExampleDirector) Release(ctx context.Context, conn *grpc.ClientConn) {
	conn.Close()
}

var _ proxy.StreamDirector = &ExampleDirector{}
