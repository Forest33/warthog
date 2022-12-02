// Package main gRPC debug server
package main

import (
	"context"
	"io"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/forest33/warthog/pkg/logger"
	testProto "github.com/forest33/warthog/testprotos"
)

const (
	addr = ":33333"
)

// Server object capable of interacting with Server
type Server struct {
	log *logger.Zerolog
}

func main() {
	s := &Server{
		log: logger.NewDefaultZerolog(),
	}

	srv := grpc.NewServer(grpc.UnaryInterceptor(s.unaryInterceptor), grpc.StreamInterceptor(s.streamInterceptor))
	testProto.RegisterTestProtoServer(srv, s)
	reflection.Register(srv)

	s.log.Debug().Msgf("server started on %s", addr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		s.log.Fatal("unable to create gRPC listener: %v", err)
	}

	if err = srv.Serve(listener); err != nil {
		s.log.Fatal("unable to start server: %v", err)
	}
}

// Unary is a Unary method handler
func (s *Server) Unary(_ context.Context, m1 *testProto.M1) (*testProto.M1, error) {
	return m1, nil
}

// TypesTest is a TypesTest method handler
func (s *Server) TypesTest(_ context.Context, t *testProto.Types) (*testProto.Types, error) {
	return t, nil
}

// LoopTest is a LoopTest method handler
func (s *Server) LoopTest(_ context.Context, t *testProto.Loop) (*testProto.Loop, error) {
	return t, nil
}

// ClientStream is a ClientStream method handler
func (s *Server) ClientStream(stream testProto.TestProto_ClientStreamServer) error {
	var (
		headers  = make([]*testProto.M3, 0, 1)
		payloads = make([]*testProto.M4, 0, 1)
	)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			s.log.Debug().Msg("caller canceled")
			err := stream.SendAndClose(&testProto.ClientStreamResponse{Header: headers, Payload: payloads})
			if err != nil {
				s.log.Error().Msgf("failed to send: %v", err)
			}
			return err
		}
		if err != nil {
			return err
		}

		header := req.GetHeader()
		payload := req.GetPayload()

		if header != nil {
			s.log.Debug().Msgf("header %+v", header)
			headers = append(headers, header)
		} else {
			s.log.Debug().Msgf("payload %+v", payload)
			payloads = append(payloads, payload)
		}
	}
}

// ServerStream is a ServerStream method handler
func (s *Server) ServerStream(req *testProto.StreamMessage, stream testProto.TestProto_ServerStreamServer) error {
	header := req.GetHeader()
	payload := req.GetPayload()

	var err error

	for stream.Context().Err() == nil {
		if header != nil {
			s.log.Debug().Msgf("send header %+v", header)
			err = stream.SendMsg(&testProto.StreamMessage{TestStream: &testProto.StreamMessage_Header{Header: header}})
		} else {
			s.log.Debug().Msgf("send payload %+v", payload)
			err = stream.SendMsg(&testProto.StreamMessage{TestStream: &testProto.StreamMessage_Payload{Payload: payload}})
		}
		if err != nil {
			s.log.Error().Msgf("failed to send message: %v", err)
			return status.Error(codes.Internal, err.Error())
		}
		time.Sleep(time.Second * 3)
	}

	s.log.Debug().Msg("caller canceled")

	return nil
}

// ClientServerStream is a ClientServerStream method handler
func (s *Server) ClientServerStream(stream testProto.TestProto_ClientServerStreamServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			s.log.Debug().Msg("caller canceled")
			return err
		}
		if err != nil {
			return err
		}

		header := req.GetHeader()
		payload := req.GetPayload()

		if header != nil {
			s.log.Debug().Msgf("send header %+v", header)
			err = stream.SendMsg(&testProto.StreamMessage{TestStream: &testProto.StreamMessage_Header{Header: header}})
		} else {
			s.log.Debug().Msgf("send payload %+v", payload)
			err = stream.SendMsg(&testProto.StreamMessage{TestStream: &testProto.StreamMessage_Payload{Payload: payload}})
		}
		if err != nil {
			s.log.Error().Msgf("failed to send message: %v", err)
			return status.Error(codes.Internal, err.Error())
		}
	}
}

func (s *Server) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ev := s.log.Debug()
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		ev = ev.Interface("metadata", md)
	}
	ev.Str("method", info.FullMethod).Msg("unary request")

	if err := grpc.SendHeader(ctx, metadata.Pairs("header-key", time.Now().String())); err != nil {
		return nil, err
	}

	if err := grpc.SetTrailer(ctx, metadata.Pairs("trailer-key", time.Now().UTC().String())); err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

func (s *Server) streamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ev := s.log.Debug()
	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		ev = ev.Interface("metadata", md)
	}
	ev.Str("method", info.FullMethod).
		Bool("is_client_stream", info.IsClientStream).
		Bool("is_server_stream", info.IsServerStream).
		Msg("stream request")

	if err := grpc.SendHeader(stream.Context(), metadata.Pairs("header-key", time.Now().String())); err != nil {
		return err
	}

	if err := grpc.SetTrailer(stream.Context(), metadata.Pairs("trailer-key", time.Now().UTC().String())); err != nil {
		return err
	}

	return handler(srv, stream)
}
