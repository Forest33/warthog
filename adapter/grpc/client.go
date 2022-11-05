package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Forest33/warthog/business/entity"
)

type Client struct {
	ctx            context.Context
	cfg            *entity.GrpcConfig
	conn           *grpc.ClientConn
	cancelQuery    context.CancelFunc
	cancelQueryMux sync.Mutex
	opts           ClientOptions
	protoPath      []string
	importPath     []string
}

func New(ctx context.Context, cfg *entity.GrpcConfig) *Client {
	return &Client{
		ctx: ctx,
		cfg: cfg,
	}
}

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
	if !c.cfg.NonBlocking {
		var cancel context.CancelFunc
		dialOptions = append(dialOptions, grpc.WithBlock())
		ctx, cancel = context.WithTimeout(c.ctx, time.Second*time.Duration(c.cfg.ConnectTimeout))
		defer cancel()
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

func (c *Client) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
