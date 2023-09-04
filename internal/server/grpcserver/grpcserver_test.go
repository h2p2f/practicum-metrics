package grpcserver

import (
	"context"
	"errors"
	"github.com/h2p2f/practicum-metrics/internal/server/grpcserver/mocks"
	pb "github.com/h2p2f/practicum-metrics/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestServer_UpdateMetric(t *testing.T) {
	tests := []struct {
		name string
		req  *pb.UpdateMetricRequest
		want *pb.UpdateMetricResponse
		err  error
	}{
		{
			name: "Positive test 1",
			req: &pb.UpdateMetricRequest{
				Metric: &pb.Metric{
					Name:    "test1",
					Type:    "gauge",
					Gauge:   1.0,
					Counter: 0,
				},
			},
			want: &pb.UpdateMetricResponse{
				Metric: &pb.Metric{
					Name:    "test1",
					Type:    "gauge",
					Gauge:   1.0,
					Counter: 0,
				},
			},
			err: nil,
		},
		{
			name: "Positive test 2",
			req: &pb.UpdateMetricRequest{
				Metric: &pb.Metric{
					Name:    "test2",
					Type:    "counter",
					Gauge:   0,
					Counter: 1,
				},
			},
			want: &pb.UpdateMetricResponse{
				Metric: &pb.Metric{
					Name:    "test2",
					Type:    "counter",
					Gauge:   0,
					Counter: 1,
				},
			},
			err: nil,
		},
		{
			name: "Negative test 1",
			req: &pb.UpdateMetricRequest{
				Metric: &pb.Metric{
					Name:    "test3",
					Type:    "gauge",
					Gauge:   -1.0,
					Counter: 0,
				},
			},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "gauge value is negative"),
		},
		{
			name: "Negative test 2",
			req: &pb.UpdateMetricRequest{
				Metric: &pb.Metric{
					Name:    "test4",
					Type:    "counter",
					Gauge:   0,
					Counter: -1,
				},
			},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "counter value is negative"),
		},
		{
			name: "Negative test 3",
			req: &pb.UpdateMetricRequest{
				Metric: &pb.Metric{
					Name:    "test5",
					Type:    "invalid_type",
					Gauge:   0,
					Counter: 0,
				},
			},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "invalid metric type"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updater := mocks.NewUpdater(t)
			logger := zap.NewNop()
			if tt.want != nil {
				switch tt.req.Metric.Type {
				case "gauge":
					updater.On("SetGauge",
						tt.req.Metric.Name, tt.req.Metric.Gauge).Return(nil)
					updater.On("GetGauge",
						tt.req.Metric.Name).Return(tt.req.Metric.Gauge, nil)
				case "counter":
					updater.On("SetCounter",
						tt.req.Metric.Name, tt.req.Metric.Counter).Return(nil)
					updater.On("GetCounter",
						tt.req.Metric.Name).Return(tt.req.Metric.Counter, nil)

				}
			}
			server := NewServer(updater, logger)
			got, err := server.UpdateMetric(context.Background(), tt.req)
			if !errors.Is(err, tt.err) {
				t.Errorf("UpdateMetric() error = %v, wantErr %v", err, tt.err)
				return
			}
			if got != nil {
				if got.Metric.Name != tt.want.Metric.Name {
					t.Errorf("UpdateMetric() got = %v, want %v", got.Metric.Name, tt.want.Metric.Name)
				}
				if got.Metric.Type != tt.want.Metric.Type {
					t.Errorf("UpdateMetric() got = %v, want %v", got.Metric.Type, tt.want.Metric.Type)
				}
				if got.Metric.Gauge != tt.want.Metric.Gauge {
					t.Errorf("UpdateMetric() got = %v, want %v", got.Metric.Gauge, tt.want.Metric.Gauge)
				}
				if got.Metric.Counter != tt.want.Metric.Counter {
					t.Errorf("UpdateMetric() got = %v, want %v", got.Metric.Counter, tt.want.Metric.Counter)
				}
			}
		})
	}
}

func TestServer_UpdateMetrics(t *testing.T) {

	tests := []struct {
		name string
		req  *pb.UpdateMetricsRequest
		want *pb.UpdateMetricsResponse
		err  error
	}{
		{
			name: "Positive test 1",
			req: &pb.UpdateMetricsRequest{
				Metrics: []*pb.Metric{
					{
						Name:    "test1",
						Type:    "gauge",
						Gauge:   1.0,
						Counter: 0,
					},
					{
						Name:    "test2",
						Type:    "counter",
						Gauge:   0,
						Counter: 1,
					},
				},
			},
			want: &pb.UpdateMetricsResponse{},
			err:  nil,
		},
		{
			name: "Negative test 1",
			req: &pb.UpdateMetricsRequest{
				Metrics: []*pb.Metric{
					{
						Name:    "test3",
						Type:    "gauge",
						Gauge:   -1.0,
						Counter: 0,
					},
				},
			},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "gauge value is negative"),
		},
		{
			name: "Negative test 2",
			req: &pb.UpdateMetricsRequest{
				Metrics: []*pb.Metric{
					{
						Name:    "test4",
						Type:    "counter",
						Gauge:   0,
						Counter: -1,
					},
				},
			},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "counter value is negative"),
		},
		{
			name: "Negative test 3",
			req: &pb.UpdateMetricsRequest{
				Metrics: []*pb.Metric{
					{
						Name:    "test5",
						Type:    "invalid_type",
						Gauge:   0,
						Counter: 0,
					},
				},
			},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "invalid metric type"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updater := mocks.NewUpdater(t)
			logger := zap.NewNop()
			if tt.want != nil {
				for _, metric := range tt.req.Metrics {
					switch metric.Type {
					case "gauge":
						updater.On("SetGauge", metric.Name, metric.Gauge).Return(nil)
					case "counter":
						updater.On("SetCounter", metric.Name, metric.Counter).Return(nil)
					}
				}
			}
			server := NewServer(updater, logger)
			_, err := server.UpdateMetrics(context.Background(), tt.req)
			if !errors.Is(err, tt.err) {
				t.Errorf("UpdateMetrics() error = %v, wantErr %v", err, tt.err)
				return
			}
		})
	}
}
