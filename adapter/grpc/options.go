// Package grpc provides basic gRPC functions.
package grpc

// ClientOptions represents Client options
type ClientOptions struct {
	noTLS              bool
	insecureSkipVerify bool
	rootCertificate    string
	clientCertificate  string
	clientKey          string
}

// ClientOpt represents Client option
type ClientOpt func(options *ClientOptions)

var defaultOptions = &ClientOptions{
	noTLS:              false,
	insecureSkipVerify: false,
	rootCertificate:    "",
	clientCertificate:  "",
	clientKey:          "",
}

// WithNoTLS returns ClientOpt which disables transport security
func WithNoTLS() ClientOpt {
	return func(options *ClientOptions) {
		options.noTLS = true
	}
}

// WithInsecure returns ClientOpt which disables server certificate chain verification and hostname
func WithInsecure() ClientOpt {
	return func(options *ClientOptions) {
		options.insecureSkipVerify = true
	}
}

// WithRootCertificate returns ClientOpt which sets server CA certificate
func WithRootCertificate(cert string) ClientOpt {
	return func(options *ClientOptions) {
		options.rootCertificate = cert
	}
}

// WithClientCertificate returns ClientOpt which sets client certificate
func WithClientCertificate(cert string) ClientOpt {
	return func(options *ClientOptions) {
		options.clientCertificate = cert
	}
}

// WithClientKey returns ClientOpt which sets client certificate private key
func WithClientKey(key string) ClientOpt {
	return func(options *ClientOptions) {
		options.clientKey = key
	}
}
