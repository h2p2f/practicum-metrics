package middlewares

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
	"testing"
)

func TestSubnetUnaryInterceptor(t *testing.T) {

	tests := []struct {
		name     string
		subnet   *net.IPNet
		ip       string
		wantErr  bool
		wantCode codes.Code
	}{
		{
			name: "Positive test 1",
			subnet: &net.IPNet{
				IP:   net.ParseIP("192.168.0.0"),
				Mask: net.CIDRMask(24, 32),
			},
			ip:       "192.168.0.1",
			wantErr:  false,
			wantCode: codes.OK,
		},
		{
			name: "Negative test 1",
			subnet: &net.IPNet{
				IP:   net.ParseIP("192.168.0.0"),
				Mask: net.CIDRMask(24, 32),
			},
			ip:       "10.0.0.1",
			wantErr:  true,
			wantCode: codes.PermissionDenied,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("X-Real-IP", tt.ip))

			req := "request"

			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return "response", nil
			}

			interceptor := SubnetUnaryInterceptor(tt.subnet)

			_, err := interceptor(ctx, req, &grpc.UnaryServerInfo{}, handler)

			if tt.wantErr && err == nil {
				t.Errorf("expected error to be not nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("expected error to be nil, got %v", err)
			}

			status, ok := status.FromError(err)
			if !ok {
				t.Errorf("expected error to be a grpc status, got %v", err)
			}
			if status.Code() != tt.wantCode {
				t.Errorf("expected error code to be %v, got %v", tt.wantCode, status.Code())
			}
		})
	}
}
