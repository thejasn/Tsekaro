package client

import (
	"context"
	"log"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

//UnaryPropagateHeaderInterceptor copy given fields from Incoming request into Outgoing request
// Empty array will make the interceptor copy all metadata in the context
func UnaryPropagateHeaderInterceptor(fields []string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			pairs := transformMapToPairs(md, fields)
			ctx = metadata.AppendToOutgoingContext(ctx, pairs...)
		}
		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}

//UnaryTimeoutInterceptor monitor the DeadlineExceeded error and log it
func UnaryTimeoutInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		return handleError(err, method, start)
	}
}

//StreamTimeoutInterceptor monitor the DeadlineExceeded error and log it
func StreamTimeoutInterceptor() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption) (grpc.ClientStream, error) {
		start := time.Now()
		stream, err := streamer(ctx, desc, cc, method, opts...)
		err = handleError(err, method, start)
		return stream, err
	}
}

func handleError(err error, method string, start time.Time) error {
	if err == nil {
		return err
	}
	statusErr, ok := status.FromError(err)
	if !ok {
		return err
	}
	if statusErr.Code() != codes.DeadlineExceeded {
		return err
	}
	log.Printf(
		"Timeout - Invoked RPC method=%s; Duration=%s; Error=%+v",
		method,
		time.Since(start), err,
	)
	return err
}

//StreamPropagateHeaderInterceptor copy given fields from Incoming request into Outgoing request
// Empty array will make the interceptor copy all metadata in the context
func StreamPropagateHeaderInterceptor(fields []string) grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption) (grpc.ClientStream, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			pairs := transformMapToPairs(md, fields)
			ctx = metadata.AppendToOutgoingContext(ctx, pairs...)
		}
		stream, err := streamer(ctx, desc, cc, method, opts...)
		return stream, err
	}
}

func transformMapToPairs(md map[string][]string, fields []string) []string {
	var kv []string
	for key, value := range md {
		if len(fields) > 0 && !contains(fields, key) {
			continue
		}
		for _, v := range value {
			kv = append(kv, key, v)
		}
	}
	return kv
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if strings.EqualFold(n, x) {
			return true
		}
	}
	return false
}
