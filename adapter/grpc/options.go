package grpc

type ClientOptions struct {
	noTLS              bool
	insecureSkipVerify bool
	rootCertificate    string
	clientCertificate  string
	clientKey          string
}

type ClientOpt func(options *ClientOptions)

var defaultOptions = &ClientOptions{
	noTLS:              false,
	insecureSkipVerify: false,
	rootCertificate:    "",
	clientCertificate:  "",
	clientKey:          "",
}

func WithNoTLS() ClientOpt {
	return func(options *ClientOptions) {
		options.noTLS = true
	}
}

func WithInsecure() ClientOpt {
	return func(options *ClientOptions) {
		options.insecureSkipVerify = true
	}
}

func WithRootCertificate(cert string) ClientOpt {
	return func(options *ClientOptions) {
		options.rootCertificate = cert
	}
}

func WithClientCertificate(cert string) ClientOpt {
	return func(options *ClientOptions) {
		options.clientCertificate = cert
	}
}

func WithClientKey(key string) ClientOpt {
	return func(options *ClientOptions) {
		options.clientKey = key
	}
}
