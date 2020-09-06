package grpc

import (
	"context"
	"io"
	"os"

	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/pkg/errors"
	"github.com/thejasn/tester/core/reflect"
	"github.com/thejasn/tester/pkg/log"
	"google.golang.org/grpc"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

var (
	exit = os.Exit
)

// To avoid confusion between program error codes and the gRPC resonse
// status codes 'Cancelled' and 'Unknown', 1 and 2 respectively,
// the response status codes emitted use an offest of 64
const statusCodeOffset = 64

type ReflectClientBuilder struct {
	ctx         context.Context
	in          io.Reader
	cc          *grpc.ClientConn
	addlHeaders multiString
	rpcHeaders  multiString
}

type multiString []string

// WithContext set the context to be used in the dial
func (b *ReflectClientBuilder) WithContext(ctx context.Context) {
	b.ctx = ctx
}

func (r *ReflectClientBuilder) WithPayload(in io.Reader) {
	r.in = in
}

func (r *ReflectClientBuilder) WithClientConn(cc *grpc.ClientConn) {
	r.cc = cc
}

func (r *ReflectClientBuilder) WithRPCHeaders(headers multiString) {
	r.rpcHeaders = headers
}

func (r *ReflectClientBuilder) WithAdditionalHeaders(headers multiString) {
	r.addlHeaders = headers
}

func (r *ReflectClientBuilder) InvokeRPC(methodName string) (string, error) {
	if r.ctx == nil {
		r.ctx = context.Background()
	}

	refClient := grpcreflect.NewClient(r.ctx, reflectpb.NewServerReflectionClient(r.cc))
	descSource := reflect.DescriptorSourceFromServer(r.ctx, refClient)

	rf, formatter, err := reflect.RequestParserAndFormatterFor(reflect.Format(reflect.FormatJSON), descSource, true, r.in)
	if err != nil {
		log.GetLogger(r.ctx).Error(errors.Wrapf(err, "Failed to construct request parser and formatter for %s", reflect.FormatJSON))
	}
	h := reflect.NewDefaultEventHandler(os.Stdout, descSource, formatter, true)

	resp, err := reflect.InvokeRPC(r.ctx, descSource, r.cc, methodName, append(r.addlHeaders, r.rpcHeaders...), h, rf.Next)
	log.GetLogger(r.ctx).Debugf("Sent: %d and received %d\n", rf.NumRequests(), h.NumResponses)

	return resp, err
	// if h.Status.Code() != codes.OK {
	// 	reflect.PrintStatus(os.Stderr, h.Status, formatter)

	// 	exit(statusCodeOffset + int(h.Status.Code()))
	// }

	// // arrange for the RPCs to be cleanly shutdown
	// reset := func() {
	// 	if refClient != nil {
	// 		refClient.Reset()
	// 		refClient = nil
	// 	}
	// 	if r.cc != nil {
	// 		r.cc.Close()
	// 		r.cc = nil
	// 	}
	// }
	// defer reset()
	// exit = func(code int) {
	// 	// since defers aren't run by os.Exit...
	// 	reset()
	// 	os.Exit(code)
	// }
}
