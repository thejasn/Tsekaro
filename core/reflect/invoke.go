package reflect

import (
	"context"
	"fmt"
	"io"
	"strings"

	protov1 "github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// InvocationEventHandler is a bag of callbacks for handling events that occur in the course
// of invoking an RPC. The handler also provides request data that is sent. The callbacks are
// generally called in the order they are listed below.
type InvocationEventHandler interface {
	// OnResolveMethod is called with a descriptor of the method that is being invoked.
	OnResolveMethod(*desc.MethodDescriptor)
	// OnSendHeaders is called with the request metadata that is being sent.
	OnSendHeaders(metadata.MD)
	// OnReceiveHeaders is called when response headers have been received.
	OnReceiveHeaders(metadata.MD)
	// OnReceiveResponse is called for each response message received.
	OnReceiveResponse(proto.Message) (string, error)
	// OnReceiveTrailers is called when response trailers and final RPC status have been received.
	OnReceiveTrailers(*status.Status, metadata.MD)
}

// RequestSupplier is a function that is called to populate messages for a gRPC operation. The
// function should populate the given message or return a non-nil error. If the supplier has no
// more messages, it should return io.EOF. When it returns io.EOF, it should not in any way
// modify the given message argument.
type RequestSupplier func(proto.Message) error

// InvokeRPC uses the given gRPC channel to invoke the given method. The given descriptor source
// is used to determine the type of method and the type of request and response message. The given
// headers are sent as request metadata. Methods on the given event handler are called as the
// invocation proceeds.
//
// The given requestData function supplies the actual data to send. It should return io.EOF when
// there is no more request data. If the method being invoked is a unary or server-streaming RPC
// (e.g. exactly one request message) and there is no request data (e.g. the first invocation of
// the function returns io.EOF), then an empty request message is sent.
//
// If the requestData function and the given event handler coordinate or share any state, they should
// be thread-safe. This is because the requestData function may be called from a different goroutine
// than the one invoking event callbacks. (This only happens for bi-directional streaming RPCs, where
// one goroutine sends request messages and another consumes the response messages).
func InvokeRPC(ctx context.Context, source DescriptorSource, ch grpcdynamic.Channel, methodName string,
	headers []string, handler InvocationEventHandler, requestData RequestSupplier) (string, error) {

	md := MetadataFromHeaders(headers)

	svc, mth := parseSymbol(methodName)
	if svc == "" || mth == "" {
		return "", fmt.Errorf("given method name %q is not in expected format: 'service/method' or 'service.method'", methodName)
	}
	dsc, err := source.FindSymbol(svc)
	if err != nil {
		if isNotFoundError(err) {
			return "", fmt.Errorf("target server does not expose service %q", svc)
		}
		return "", fmt.Errorf("failed to query for service descriptor %q: %v", svc, err)
	}
	sd, ok := dsc.(*desc.ServiceDescriptor)
	if !ok {
		return "", fmt.Errorf("target server does not expose service %q", svc)
	}
	mtd := sd.FindMethodByName(mth)
	if mtd == nil {
		return "", fmt.Errorf("service %q does not include a method named %q", svc, mth)
	}

	//handler.OnResolveMethod(mtd)

	// we also download any applicable extensions so we can provide full support for parsing user-provided data
	var ext dynamic.ExtensionRegistry
	alreadyFetched := map[string]bool{}
	if err = fetchAllExtensions(source, &ext, mtd.GetInputType(), alreadyFetched); err != nil {
		return "", fmt.Errorf("error resolving server extensions for message %s: %v", mtd.GetInputType().GetFullyQualifiedName(), err)
	}
	if err = fetchAllExtensions(source, &ext, mtd.GetOutputType(), alreadyFetched); err != nil {
		return "", fmt.Errorf("error resolving server extensions for message %s: %v", mtd.GetOutputType().GetFullyQualifiedName(), err)
	}

	msgFactory := dynamic.NewMessageFactoryWithExtensionRegistry(&ext)
	req := msgFactory.NewMessage(mtd.GetInputType())

	//handler.OnSendHeaders(md)
	ctx = metadata.NewOutgoingContext(ctx, md)

	stub := grpcdynamic.NewStubWithMessageFactory(ch, msgFactory)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return invokeUnary(ctx, stub, mtd, handler, requestData, protov1.MessageV2(req))
}

func invokeUnary(ctx context.Context, stub grpcdynamic.Stub, md *desc.MethodDescriptor, handler InvocationEventHandler,
	requestData RequestSupplier, req proto.Message) (string, error) {

	err := requestData(req)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("error getting request data: %v", err)
	}
	if err != io.EOF {
		// verify there is no second message, which is a usage error
		err := requestData(req)
		if err == nil {
			return "", fmt.Errorf("method %q is a unary RPC, but request data contained more than 1 message", md.GetFullyQualifiedName())
		} else if err != io.EOF {
			return "", fmt.Errorf("error getting request data: %v", err)
		}
	}

	// Now we can actually invoke the RPC!
	var respHeaders metadata.MD
	var respTrailers metadata.MD
	resp, err := stub.InvokeRpc(ctx, md, protov1.MessageV1(req), grpc.Trailer(&respTrailers), grpc.Header(&respHeaders))

	stat, ok := status.FromError(err)
	if !ok {
		// Error codes sent from the server will get printed differently below.
		// So just bail for other kinds of errors here.
		return "", fmt.Errorf("grpc call for %q failed: %v", md.GetFullyQualifiedName(), err)
	}

	var r string
	if stat.Code() == codes.OK {
		return handler.OnReceiveResponse(protov1.MessageV2(resp))
	}
	return r, nil
}

type notFoundError string

func notFound(kind, name string) error {
	return notFoundError(fmt.Sprintf("%s not found: %s", kind, name))
}

func (e notFoundError) Error() string {
	return string(e)
}

func isNotFoundError(err error) bool {
	if grpcreflect.IsElementNotFoundError(err) {
		return true
	}
	_, ok := err.(notFoundError)
	return ok
}

func parseSymbol(svcAndMethod string) (string, string) {
	pos := strings.LastIndex(svcAndMethod, "/")
	if pos < 0 {
		pos = strings.LastIndex(svcAndMethod, ".")
		if pos < 0 {
			return "", ""
		}
	}
	return svcAndMethod[:pos], svcAndMethod[pos+1:]
}
