package grpc

import (
	"context"
	"strings"
	"time"

	"github.com/thejasn/tester/core/client"

	"github.com/pkg/errors"
	"github.com/thejasn/tester/pkg/log"
	"google.golang.org/grpc"
)

type Config struct {
	ctx     context.Context
	rc      ReflectClientBuilder
	key     string
	host    string
	port    string
	method  string
	request string
}

func NewConfig(ctx context.Context, key, host, port string) *Config {
	return &Config{
		ctx:  ctx,
		key:  key,
		host: host,
		port: port,
	}
}

func (p Config) GetIdentifier() string {
	return p.key
}

func WithRequest(req string) client.RunnerOpts {
	return func(p client.Runner) {
		p.(*Config).request = req
	}
}

func WithMethod(method string) client.RunnerOpts {
	return func(p client.Runner) {
		p.(*Config).method = method
	}
}

func (p *Config) Clear() {
	p.method = ""
	p.request = ""
}

func (p *Config) Build() error {
	dial := func() *grpc.ClientConn {
		clientBuilder := GrpcClientBuilder{}
		dialTime := 10 * time.Second
		ctx, cancel := context.WithTimeout(p.ctx, dialTime)
		defer cancel()
		clientBuilder.WithContext(ctx)
		cc, err := clientBuilder.GetConn(p.host, p.port)
		if err != nil {
			log.G(ctx).Fatal(errors.Wrapf(err, "Failed to dial target host %q and port %q", p.host, p.port))
		}
		return cc
	}
	p.rc = ReflectClientBuilder{}
	p.rc.WithClientConn(dial())
	p.rc.WithContext(p.ctx)
	p.rc.WithPayload(strings.NewReader(p.request))
	return nil
}

func (p *Config) Invoke() (string, error) {
	r, err := p.rc.InvokeRPC(p.method)
	if err != nil {
		return "", errors.Wrapf(err, "Error invoking method %q", p.method)
	}
	return r, nil
}
