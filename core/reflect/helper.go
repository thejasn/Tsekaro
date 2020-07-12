package reflect

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoprint"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// MetadataFromHeaders converts a list of header strings (each string in
// "Header-Name: Header-Value" form) into metadata. If a string has a header
// name without a value (e.g. does not contain a colon), the value is assumed
// to be blank. Binary headers (those whose names end in "-bin") should be
// base64-encoded. But if they cannot be base64-decoded, they will be assumed to
// be in raw form and used as is.
func MetadataFromHeaders(headers []string) metadata.MD {
	md := make(metadata.MD)
	for _, part := range headers {
		if part != "" {
			pieces := strings.SplitN(part, ":", 2)
			if len(pieces) == 1 {
				pieces = append(pieces, "") // if no value was specified, just make it "" (maybe the header value doesn't matter)
			}
			headerName := strings.ToLower(strings.TrimSpace(pieces[0]))
			val := strings.TrimSpace(pieces[1])
			if strings.HasSuffix(headerName, "-bin") {
				if v, err := decode(val); err == nil {
					val = v
				}
			}
			md[headerName] = append(md[headerName], val)
		}
	}
	return md
}

var base64Codecs = []*base64.Encoding{base64.StdEncoding, base64.URLEncoding, base64.RawStdEncoding, base64.RawURLEncoding}

func decode(val string) (string, error) {
	var firstErr error
	var b []byte
	// we are lenient and can accept any of the flavors of base64 encoding
	for _, d := range base64Codecs {
		var err error
		b, err = d.DecodeString(val)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		return string(b), nil
	}
	return "", firstErr
}

// fetchAllExtensions recursively fetches from the server extensions for the given message type as well as
// for all message types of nested fields. The extensions are added to the given dynamic registry of extensions
// so that all server-known extensions can be correctly parsed..
func fetchAllExtensions(source DescriptorSource, ext *dynamic.ExtensionRegistry, md *desc.MessageDescriptor, alreadyFetched map[string]bool) error {
	msgTypeName := md.GetFullyQualifiedName()
	if alreadyFetched[msgTypeName] {
		return nil
	}
	alreadyFetched[msgTypeName] = true
	if len(md.GetExtensionRanges()) > 0 {
		fds, err := source.AllExtensionsForType(msgTypeName)
		if err != nil {
			return fmt.Errorf("failed to query for extensions of type %s: %v", msgTypeName, err)
		}
		for _, fd := range fds {
			if err := ext.AddExtension(fd); err != nil {
				return fmt.Errorf("could not register extension %s of type %s: %v", fd.GetFullyQualifiedName(), msgTypeName, err)
			}
		}
	}
	// recursively fetch extensions for the types of any message fields
	for _, fd := range md.GetFields() {
		if fd.GetMessageType() != nil {
			err := fetchAllExtensions(source, ext, fd.GetMessageType(), alreadyFetched)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// MetadataToString returns a string representation of the given metadata, for
// displaying to users.
func MetadataToString(md metadata.MD) string {
	if len(md) == 0 {
		return "(empty)"
	}

	keys := make([]string, 0, len(md))
	for k := range md {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b bytes.Buffer
	first := true
	for _, k := range keys {
		vs := md[k]
		for _, v := range vs {
			if first {
				first = false
			} else {
				b.WriteString("\n")
			}
			b.WriteString(k)
			b.WriteString(": ")
			if strings.HasSuffix(k, "-bin") {
				v = base64.StdEncoding.EncodeToString([]byte(v))
			}
			b.WriteString(v)
		}
	}
	return b.String()
}

var printer = &protoprint.Printer{
	Compact:                  true,
	OmitComments:             protoprint.CommentsNonDoc,
	SortElements:             true,
	ForceFullyQualifiedNames: true,
}

// GetDescriptorText returns a string representation of the given descriptor.
// This returns a snippet of proto source that describes the given element.
func GetDescriptorText(dsc desc.Descriptor, _ DescriptorSource) (string, error) {
	// Note: DescriptorSource is not used, but remains an argument for backwards
	// compatibility with previous implementation.
	txt, err := printer.PrintProtoToString(dsc)
	if err != nil {
		return "", err
	}
	// callers don't expect trailing newlines
	if txt[len(txt)-1] == '\n' {
		txt = txt[:len(txt)-1]
	}
	return txt, nil
}

// errSignalingCreds is a wrapper around a TransportCredentials value, but
// it will use the writeResult function to notify on error.
type errSignalingCreds struct {
	credentials.TransportCredentials
	writeResult func(res interface{})
}

func (c *errSignalingCreds) ClientHandshake(ctx context.Context, addr string, rawConn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	conn, auth, err := c.TransportCredentials.ClientHandshake(ctx, addr, rawConn)
	if err != nil {
		c.writeResult(err)
	}
	return conn, auth, err
}

// BlockingDial is a helper method to dial the given address, using optional TLS credentials,
// and blocking until the returned connection is ready. If the given credentials are nil, the
// connection will be insecure (plain-text).
func BlockingDial(ctx context.Context, network, address string, creds credentials.TransportCredentials, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	// grpc.Dial doesn't provide any information on permanent connection errors (like
	// TLS handshake failures). So in order to provide good error messages, we need a
	// custom dialer that can provide that info. That means we manage the TLS handshake.
	result := make(chan interface{}, 1)

	writeResult := func(res interface{}) {
		// non-blocking write: we only need the first result
		select {
		case result <- res:
		default:
		}
	}

	// custom credentials and dialer will notify on error via the
	// writeResult function
	if creds != nil {
		creds = &errSignalingCreds{
			TransportCredentials: creds,
			writeResult:          writeResult,
		}
	}
	dialer := func(ctx context.Context, address string) (net.Conn, error) {
		// NB: We *could* handle the TLS handshake ourselves, in the custom
		// dialer (instead of customizing both the dialer and the credentials).
		// But that requires using WithInsecure dial option (so that the gRPC
		// library doesn't *also* try to do a handshake). And that would mean
		// that the library would send the wrong ":scheme" metaheader to
		// servers: it would send "http" instead of "https" because it is
		// unaware that TLS is actually in use.
		conn, err := (&net.Dialer{}).DialContext(ctx, network, address)
		if err != nil {
			writeResult(err)
		}
		return conn, err
	}

	// Even with grpc.FailOnNonTempDialError, this call will usually timeout in
	// the face of TLS handshake errors. So we can't rely on grpc.WithBlock() to
	// know when we're done. So we run it in a goroutine and then use result
	// channel to either get the connection or fail-fast.
	go func() {
		opts = append(opts,
			grpc.WithBlock(),
			grpc.FailOnNonTempDialError(true),
			grpc.WithContextDialer(dialer),
		)
		if creds == nil {
			opts = append(opts, grpc.WithInsecure())
		} else {
			opts = append(opts, grpc.WithTransportCredentials(creds))
		}
		conn, err := grpc.DialContext(ctx, address, opts...)
		var res interface{}
		if err != nil {
			res = err
		} else {
			res = conn
		}
		writeResult(res)
	}()

	select {
	case res := <-result:
		if conn, ok := res.(*grpc.ClientConn); ok {
			return conn, nil
		}
		return nil, res.(error)
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
