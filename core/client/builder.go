package client

import (
	"context"
	"fmt"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	log "github.com/sirupsen/logrus"
	"github.com/thejasn/tester/core/reflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

//GrpcClientConnBuilder is a builder to create GRPC connection to the GRPC Server
type GrpcClientConnBuilder interface {
	WithContext(ctx context.Context)
	WithOptions(opts ...grpc.DialOption)
	WithInsecure()
	WithTLS(credentials.TransportCredentials)
	WithUnaryInterceptors(interceptors []grpc.UnaryClientInterceptor)
	WithStreamInterceptors(interceptors []grpc.StreamClientInterceptor)
	WithKeepAliveParams(params keepalive.ClientParameters)
	GetConn(addr string, port string) (*grpc.ClientConn, error)
}

//GRPC client builder
type GrpcClientBuilder struct {
	creds   credentials.TransportCredentials
	options []grpc.DialOption
	ctx     context.Context
}

// WithContext set the context to be used in the dial
func (b *GrpcClientBuilder) WithContext(ctx context.Context) {
	b.ctx = ctx
}

// WithOptions set dial options
func (b *GrpcClientBuilder) WithOptions(opts ...grpc.DialOption) {
	b.options = append(b.options, opts...)
}

// WithInsecure set the connection as insecure
func (b *GrpcClientBuilder) WithInsecure() {
	b.options = append(b.options, grpc.WithInsecure())
}

//WithTLS creates a connection using the provided tls credentials
func (b *GrpcClientBuilder) WithTLS(creds credentials.TransportCredentials) {
	b.creds = creds
}

// WithKeepAliveParams set the keep alive params
// ClientParameters is used to set keepalive parameters on the client-side.
// These configure how the client will actively probe to notice when a
// connection is broken and send pings so intermediaries will be aware of the
// liveness of the connection. Make sure these parameters are set in
// coordination with the keepalive policy on the server, as incompatible
// settings can result in closing of connection.
func (b *GrpcClientBuilder) WithKeepAliveParams(params keepalive.ClientParameters) {
	keepAlive := grpc.WithKeepaliveParams(params)
	b.options = append(b.options, keepAlive)
}

// WithUnaryInterceptors set a list of interceptors to the Grpc client for unary connection
// By default, gRPC doesn't allow one to have more than one interceptor either on the client nor on the server side.
// By using `grpc_middleware` we are able to provides convenient method to add a list of interceptors
func (b *GrpcClientBuilder) WithUnaryInterceptors(interceptors []grpc.UnaryClientInterceptor) {
	b.options = append(b.options, grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(interceptors...)))
}

// WithUnaryInterceptors set a list of interceptors to the Grpc client for stream connection
// By default, gRPC doesn't allow one to have more than one interceptor either on the client nor on the server side.
// By using `grpc_middleware` we are able to provides convenient method to add a list of interceptors
func (b *GrpcClientBuilder) WithStreamInterceptors(interceptors []grpc.StreamClientInterceptor) {
	b.options = append(b.options, grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(interceptors...)))
}

// GetConn returns the client connection to the server
func (b *GrpcClientBuilder) GetConn(addr string, port string) (*grpc.ClientConn, error) {
	if addr == "" || port == "" {
		return nil, fmt.Errorf("target connection parameter missing. address = %s, port = %s", addr, port)
	}
	target := strings.Join([]string{addr, port}, ":")
	log.Debugf("Target to connect = %s", target)
	ctx := b.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	cc, err := reflect.BlockingDial(ctx, "tcp", target, b.creds)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to client. address = %s, port = %s. error = %+v", addr, port, err)
	}
	return cc, nil
}
