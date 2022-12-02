// Package grpc provides basic gRPC functions.
package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"sync"
	"time"

	"github.com/jhump/protoreflect/dynamic"

	"github.com/forest33/warthog/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/forest33/warthog/business/entity"
)

const (
	requestChanCapacity  = 10
	responseChanCapacity = 10
)

// Client object capable of interacting with Client
type Client struct {
	ctx              context.Context
	cfg              *entity.Settings
	log              *logger.Zerolog
	conn             *grpc.ClientConn
	queryCtx         context.Context
	queryCancel      context.CancelFunc
	queryStartTime   time.Time
	cancelQueryMux   sync.Mutex
	requestCh        chan *dynamic.Message
	responseCh       chan *entity.QueryResponse
	closeStreamCh    chan struct{}
	sentMessages     uint
	receivedMessaged uint
	opts             ClientOptions
	protoPath        []string
	importPath       []string
}

// New creates a new Client
func New(ctx context.Context, log *logger.Zerolog) *Client {
	return &Client{
		ctx:           ctx,
		log:           log,
		responseCh:    make(chan *entity.QueryResponse, responseChanCapacity),
		closeStreamCh: make(chan struct{}, 1),
	}
}

// SetSettings sets application settings
func (c *Client) SetSettings(cfg *entity.Settings) {
	c.cfg = cfg
}

// Connect connecting to gRPC server
func (c *Client) Connect(addr string, opts ...ClientOpt) error {
	if defaultOptions != nil {
		c.opts = *defaultOptions
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&c.opts)
	}

	dialOptions, err := c.getDialOptions()
	if err != nil {
		return err
	}

	ctx := c.ctx
	if !*c.cfg.NonBlockingConnection {
		var cancel context.CancelFunc
		dialOptions = append(dialOptions, grpc.WithBlock())
		if *c.cfg.ConnectTimeout > 0 {
			ctx, cancel = context.WithTimeout(c.ctx, time.Second*time.Duration(*c.cfg.ConnectTimeout))
			defer cancel()
		}
	}

	c.conn, err = grpc.DialContext(ctx, addr, dialOptions...)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) getDialOptions() ([]grpc.DialOption, error) {
	if c.opts.noTLS {
		return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}, nil
	}

	creds, err := c.loadTLSCredentials()
	if err != nil {
		return nil, err
	}

	return []grpc.DialOption{grpc.WithTransportCredentials(creds)}, nil
}

func (c *Client) loadTLSCredentials() (credentials.TransportCredentials, error) {
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM([]byte(c.opts.rootCertificate)) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	clientCert, err := tls.X509KeyPair([]byte(c.opts.clientCertificate), []byte(c.opts.clientKey))
	if err != nil {
		return nil, err
	}

	cfg := &tls.Config{
		Certificates:       []tls.Certificate{clientCert},
		RootCAs:            pool,
		InsecureSkipVerify: c.opts.insecureSkipVerify,
	}

	return credentials.NewTLS(cfg), nil
}

// Close closes connection to gRPC server
func (c *Client) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			c.log.Error().Msgf("failed to close connection: %v", err)
		}
	}
}
