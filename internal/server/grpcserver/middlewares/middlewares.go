// Package middlewares implement wrappers for grpc server with logging and checking IP.
package middlewares

import (
	"context"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
)

// codeToLevel function for converting grpc codes to zap levels
// if the code is OK, then the level is DebugLevel
// otherwise, use the DefaultCodeToLevel function from grpc_zap
func codeToLevel(code codes.Code) zapcore.Level {
	if code == codes.OK {
		return zap.DebugLevel
	}
	return grpc_zap.DefaultCodeToLevel(code)
}

// WithLogging adds logging middleware to grpc server opts
func WithLogging(logger *zap.Logger, opts []grpc.ServerOption) []grpc.ServerOption {
	opts = append(
		opts,
		grpc.UnaryInterceptor(grpc_zap.UnaryServerInterceptor(logger, grpc_zap.WithLevels(codeToLevel))))
	return opts
}

// WithCheckingIP adds checking IP middleware to grpc server opts
func WithCheckingIP(subnet *net.IPNet, opts []grpc.ServerOption) []grpc.ServerOption {
	opts = append(opts, grpc.UnaryInterceptor(SubnetUnaryInterceptor(subnet)))
	return opts
}

func WithChekingIPAndLogging(logger *zap.Logger, subnet *net.IPNet, opts []grpc.ServerOption) []grpc.ServerOption {
	opts = append(
		opts,
		grpc.ChainUnaryInterceptor(SubnetUnaryInterceptor(subnet),
			grpc_zap.UnaryServerInterceptor(logger, grpc_zap.WithLevels(codeToLevel))))
	return opts
}

// SubnetUnaryInterceptor checks if the client IP is in the subnet
func SubnetUnaryInterceptor(subnet *net.IPNet) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		// Checking if the client IP is in the subnet
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "metadata not found")
		}
		ips := md.Get("X-Real-IP")
		if len(ips) == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "X-Real-IP not found in metadata")
		}
		ip := ips[0]
		remoteIP := net.ParseIP(ip)
		if !subnet.Contains(remoteIP) {
			return nil, status.Errorf(codes.PermissionDenied, "client IP is not allowed")
		}
		return handler(ctx, req)
	}
}
