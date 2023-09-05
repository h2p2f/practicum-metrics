// Package middlewares implement wrappers for grpc server for checking IP at server side.
package middlewares

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// IPInjectorUnaryClientInterceptor adds IP to metadata
func IPInjectorUnaryClientInterceptor(addr *net.IP) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var md metadata.MD
		// check if metadata exists
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			//if not, create new metadata
			md = metadata.New(nil)
		}
		// add IP to metadata and put it to context
		md.Set("X-Real-IP", addr.String())
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
