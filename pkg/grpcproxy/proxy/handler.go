// Copyright 2017-2018 Valient Gough
// Copyright 2017 Michal Witkowski
// All Rights Reserved.
// See LICENSE for licensing terms.

package proxy

import (
	"context"
	"io"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

var (
	clientStreamDescForProxying = &grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}
)

// RegisterService sets up a proxy handler for a particular gRPC service and method.
// The behavior is the same as if you were registering a handler method, e.g. from a codegenerated pb.go file.
//
// This can *only* be used if the `server` also uses proxy.CodecForServer() ServerOption.
func RegisterService(server *grpc.Server, director StreamDirector, serviceName string, methodNames ...string) {
	streamer := &handler{director}
	fakeDesc := &grpc.ServiceDesc{
		ServiceName: serviceName,
		HandlerType: (*interface{})(nil),
	}
	for _, m := range methodNames {
		streamDesc := grpc.StreamDesc{
			StreamName:    m,
			Handler:       streamer.handler,
			ServerStreams: true,
			ClientStreams: true,
		}
		fakeDesc.Streams = append(fakeDesc.Streams, streamDesc)
	}
	server.RegisterService(fakeDesc, streamer)
}

// TransparentHandler returns a handler that attempts to proxy all requests that are not registered in the server.
// The indented use here is as a transparent proxy, where the server doesn't know about the services implemented by the
// backends. It should be used as a `grpc.UnknownServiceHandler`.
//
// This can *only* be used if the `server` also uses proxy.CodecForServer() ServerOption.
func TransparentHandler(director StreamDirector) grpc.StreamHandler {
	streamer := &handler{director}
	return streamer.handler
}

type handler struct {
	director StreamDirector
}

// handler is where the real magic of proxying happens.
// It is invoked like any gRPC server stream and uses the gRPC server framing to get and receive bytes from the wire,
// forwarding it to a ClientStream established against the relevant ClientConn.
func (s *handler) handler(srv interface{}, serverStream grpc.ServerStream) error {
	serverCtx := serverStream.Context()
	ss := grpc.ServerTransportStreamFromContext(serverCtx)
	fullMethodName := ss.Method()
	outCtx, backendConn, err := s.director.Connect(serverCtx, fullMethodName)
	if err != nil {
		return err
	}
	defer s.director.Release(outCtx, backendConn)

	clientCtx, clientCancel := context.WithCancel(outCtx)
	defer clientCancel()
	if _, ok := metadata.FromOutgoingContext(outCtx); !ok {
		clientCtx = copyMetadata(clientCtx, outCtx)
	}

	var copts []grpc.CallOption
	if md, ok := metadata.FromOutgoingContext(clientCtx); ok {
		if vals, ok := md["content-type"]; ok && len(vals) > 0 {
			if contentSubtype, ok := contentSubtype(vals[0]); ok {
				copts = append(copts, grpc.CallContentSubtype(contentSubtype))
			}
		}
	}

	clientStream, err := grpc.NewClientStream(clientCtx, clientStreamDescForProxying, backendConn, fullMethodName, copts...)
	if err != nil {
		return err
	}

	err = biDirCopy(serverStream, clientStream)
	if err == io.EOF {
		return nil
	}
	return err
}

// copyMetadata takes the new client (outgoing) context, a server (incoming)
// context, and returns a new outgoing context which contains all the incoming
// metadata.
//
// An additional X-Forwarded-For metadata entry is added or appended to with
// the peer address from the server context. See https://en.wikipedia.org/wiki/X-Forwarded-For.
func copyMetadata(ctx context.Context, serverCtx context.Context) context.Context {
	source := "unknown"
	if peer, ok := peer.FromContext(serverCtx); ok && peer.Addr != nil {
		source = peer.Addr.String()
	}
	forwardMD := metadata.Pairs("X-Forwarded-For", source)

	md, ok := metadata.FromIncomingContext(serverCtx)
	if ok {
		return metadata.NewOutgoingContext(ctx, metadata.Join(md, forwardMD))
	}
	return metadata.NewOutgoingContext(ctx, forwardMD)
}

// contentSubtype returns the content-subtype for the given content-type.  The
// given content-type must be a valid content-type that starts with
// "application/grpc". A content-subtype will follow "application/grpc" after a
// "+" or ";". See
// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#requests for
// more details.
//
// If contentType is not a valid content-type for gRPC, the boolean
// will be false, otherwise true. If content-type == "application/grpc",
// "application/grpc+", or "application/grpc;", the boolean will be true,
// but no content-subtype will be returned.
//
// contentType is assumed to be lowercase already.
func contentSubtype(contentType string) (string, bool) {
	const baseContentType = "application/grpc"
	if contentType == baseContentType {
		return "", true
	}
	if !strings.HasPrefix(contentType, baseContentType) {
		return "", false
	}
	// guaranteed since != baseContentType and has baseContentType prefix
	switch contentType[len(baseContentType)] {
	case '+', ';':
		// this will return true for "application/grpc+" or "application/grpc;"
		// which the previous validContentType function tested to be valid, so we
		// just say that no content-subtype is specified in this case
		return contentType[len(baseContentType)+1:], true
	default:
		return "", false
	}
}
