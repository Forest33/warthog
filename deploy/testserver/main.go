// nolint:gosec
// Package main gRPC debug server
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"net"
	"os"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/genproto/protobuf/ptype"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/apipb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/sourcecontextpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/forest33/warthog/pkg/logger"
	testProto "github.com/forest33/warthog/testprotos"
)

const (
	addr    = ":33333"
	withTLS = false
)

// Server object capable of interacting with Server.
type Server struct {
	log *logger.Zerolog
}

func main() {
	s := &Server{
		log: logger.NewDefaultZerolog(),
	}

	opts := []grpc.ServerOption{grpc.UnaryInterceptor(s.unaryInterceptor), grpc.StreamInterceptor(s.streamInterceptor)}

	if withTLS {
		tlsCredentials, err := loadTLSCredentials()
		if err != nil {
			s.log.Fatal("cannot load TLS credentials: ", err)
		}
		opts = append(opts, grpc.Creds(tlsCredentials))
	}

	srv := grpc.NewServer(opts...)
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

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	pemClientCA, err := os.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, errors.New("failed to add client CA's certificate")
	}

	serverCert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}

// Unary is a Unary method handler.
func (s *Server) Unary(_ context.Context, m1 *testProto.M1) (*testProto.M1, error) {
	return m1, nil
}

// CreateUser is a CreateUser method handler.
func (s *Server) CreateUser(_ context.Context, u *testProto.User) (*testProto.User, error) {
	return u, nil
}

// TypesTest is a TypesTest method handler.
func (s *Server) TypesTest(_ context.Context, t *testProto.Types) (*testProto.Types, error) {
	return t, nil
}

// LoopTest is a LoopTest method handler.
func (s *Server) LoopTest(_ context.Context, t *testProto.Loop) (*testProto.Loop, error) {
	return t, nil
}

// ClientStream is a ClientStream method handler.
func (s *Server) ClientStream(stream testProto.TestProto_ClientStreamServer) error {
	var (
		headers  = make([]*testProto.M3, 0, 1)
		payloads = make([]*testProto.M4, 0, 1)
	)

	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
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

// ServerStream is a ServerStream method handler.
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

// ClientServerStream is a ClientServerStream method handler.
func (s *Server) ClientServerStream(stream testProto.TestProto_ClientServerStreamServer) error {
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
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

// AuthBasic is a AuthBasic method handler.
func (s *Server) AuthBasic(_ context.Context, m1 *testProto.M1) (*testProto.M1, error) {
	return m1, nil
}

// AuthBearer is a AuthBearer method handler.
func (s *Server) AuthBearer(_ context.Context, m1 *testProto.M1) (*testProto.M1, error) {
	return m1, nil
}

// AuthJWT is a AuthJWT method handler.
func (s *Server) AuthJWT(_ context.Context, m1 *testProto.M1) (*testProto.M1, error) {
	return m1, nil
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

	var err error
	switch info.FullMethod {
	case "/test.proto.v1.test_proto/AuthBasic":
		err = ensureValidBasicCredentials(ctx)
	case "/test.proto.v1.test_proto/AuthBearer":
		err = ensureValidToken(ctx)
	case "/test.proto.v1.test_proto/AuthJWT":
		err = ensureValidJWT(ctx)
	}

	if err != nil {
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

// Any is a Any method handler.
func (s *Server) Any(_ context.Context, r *anypb.Any) (*anypb.Any, error) {
	return r, nil
}

// Api is a Api method handler.
func (s *Server) Api(_ context.Context, r *apipb.Api) (*apipb.Api, error) {
	return r, nil
}

// BoolValue is a BoolValue method handler.
func (s *Server) BoolValue(_ context.Context, r *wrappers.BoolValue) (*wrappers.BoolValue, error) {
	return r, nil
}

// BytesValue is a BytesValue method handler.
func (s *Server) BytesValue(_ context.Context, r *wrappers.BytesValue) (*wrappers.BytesValue, error) {
	return r, nil
}

// DoubleValue is a DoubleValue method handler.
func (s *Server) DoubleValue(_ context.Context, r *wrappers.DoubleValue) (*wrappers.DoubleValue, error) {
	return r, nil
}

// Duration is a Duration method handler.
func (s *Server) Duration(_ context.Context, r *durationpb.Duration) (*durationpb.Duration, error) {
	return r, nil
}

// Empty is a Empty method handler.
func (s *Server) Empty(_ context.Context, _ *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

// Enum is a Enum method handler.
func (s *Server) Enum(_ context.Context, r *ptype.Enum) (*ptype.Enum, error) {
	return r, nil
}

// EnumValue is a EnumValue method handler.
func (s *Server) EnumValue(_ context.Context, r *ptype.EnumValue) (*ptype.EnumValue, error) {
	return r, nil
}

// Field is a Field method handler.
func (s *Server) Field(_ context.Context, r *ptype.Field) (*ptype.Field, error) {
	return r, nil
}

// FieldMask is a FieldMask method handler.
func (s *Server) FieldMask(_ context.Context, r *fieldmaskpb.FieldMask) (*fieldmaskpb.FieldMask, error) {
	return r, nil
}

// FloatValue is a FloatValue method handler.
func (s *Server) FloatValue(_ context.Context, r *wrappers.FloatValue) (*wrappers.FloatValue, error) {
	return r, nil
}

// Int32Value is a Int32Value method handler.
func (s *Server) Int32Value(_ context.Context, r *wrappers.Int32Value) (*wrappers.Int32Value, error) {
	return r, nil
}

// Int64Value is a Int64Value method handler.
func (s *Server) Int64Value(_ context.Context, r *wrappers.Int64Value) (*wrappers.Int64Value, error) {
	return r, nil
}

// ListValue is a ListValue method handler.
func (s *Server) ListValue(_ context.Context, r *structpb.ListValue) (*structpb.ListValue, error) {
	return r, nil
}

// Method is a Method method handler.
func (s *Server) Method(_ context.Context, r *apipb.Method) (*apipb.Method, error) {
	return r, nil
}

// Mixin is a Mixin method handler.
func (s *Server) Mixin(_ context.Context, r *apipb.Mixin) (*apipb.Mixin, error) {
	return r, nil
}

// Option is a Option method handler.
func (s *Server) Option(_ context.Context, r *ptype.Option) (*ptype.Option, error) {
	return r, nil
}

// SourceContext is a SourceContext method handler.
func (s *Server) SourceContext(_ context.Context, r *sourcecontextpb.SourceContext) (*sourcecontextpb.SourceContext, error) {
	return r, nil
}

// StringValue is a StringValue method handler.
func (s *Server) StringValue(_ context.Context, r *wrappers.StringValue) (*wrappers.StringValue, error) {
	return r, nil
}

// Struct is a Struct method handler.
func (s *Server) Struct(_ context.Context, r *structpb.Struct) (*structpb.Struct, error) {
	return r, nil
}

// Timestamp is a Timestamp method handler.
func (s *Server) Timestamp(_ context.Context, r *timestamppb.Timestamp) (*timestamppb.Timestamp, error) {
	return r, nil
}

// Type is a Type method handler.
func (s *Server) Type(_ context.Context, r *ptype.Type) (*ptype.Type, error) {
	return r, nil
}

// UInt32Value is a UInt32Value method handler.
func (s *Server) UInt32Value(_ context.Context, r *wrappers.UInt32Value) (*wrappers.UInt32Value, error) {
	return r, nil
}

// UInt64Value is a UInt64Value method handler.
func (s *Server) UInt64Value(_ context.Context, r *wrappers.UInt64Value) (*wrappers.UInt64Value, error) {
	return r, nil
}

// Value is a Value method handler.
func (s *Server) Value(_ context.Context, r *structpb.Value) (*structpb.Value, error) {
	return r, nil
}
