// Package: grpcserver implements the logic of receiving metrics from the agent.

package grpcserver

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/h2p2f/practicum-metrics/proto"
)

// Updater interface for updating metrics
//
//go:generate mockery --name Updater --output ./mocks --filename mocks_updater.go
type Updater interface {
	SetGauge(name string, value float64)
	SetCounter(name string, value int64)
	GetCounter(name string) (value int64, err error)
	GetGauge(name string) (value float64, err error)
}

// Server structure for implementing the server
type Server struct {
	pb.UnimplementedMetricsServiceServer
	db     Updater
	logger *zap.Logger
}

// NewServer creates a new server
func NewServer(db Updater, logger *zap.Logger) *Server {
	return &Server{db: db, logger: logger}
}

// UpdateMetric - update metric from agent
func (s *Server) UpdateMetric(
	ctx context.Context,
	req *pb.UpdateMetricRequest,
) (*pb.UpdateMetricResponse, error) {
	// Processing and validation of the received data
	var response pb.UpdateMetricResponse
	var err error
	s.logger.Debug(
		"request from client:",
		zap.String("metric", req.Metric.Name),
		zap.String("type", req.Metric.Type),
		zap.Float64("gauge", req.Metric.Gauge),
		zap.Int64("counter", req.Metric.Counter))
	switch req.Metric.Type {
	case "gauge":
		if req.Metric.Gauge < 0 {
			return nil, status.Error(codes.InvalidArgument, "gauge value is negative")
		} else {
			s.db.SetGauge(req.Metric.Name, req.Metric.Gauge)

			response.Metric = req.Metric
			response.Metric.Gauge, err = s.db.GetGauge(req.Metric.Name)
			if err != nil {
				return nil, status.Error(codes.Internal, "can't get gauge value")
			}
		}
	case "counter":
		if req.Metric.Counter < 0 {
			return nil, status.Error(codes.InvalidArgument, "counter value is negative")
		} else {
			s.db.SetCounter(req.Metric.Name, req.Metric.Counter)
			response.Metric = req.Metric
			response.Metric.Counter, err = s.db.GetCounter(req.Metric.Name)
			if err != nil {
				return nil, status.Error(codes.Internal, "can't get counter value")
			}
		}
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid metric type")
	}
	return &response, nil
}

// UpdateMetrics - update metrics from agent in batch
func (s *Server) UpdateMetrics(
	ctx context.Context,
	req *pb.UpdateMetricsRequest,
) (*pb.UpdateMetricsResponse, error) {
	// Processing and validation of the received data
	var response pb.UpdateMetricsResponse
	s.logger.Info(
		"request from client:",
		zap.Int("number of metrics", len(req.Metrics)))
	for _, metric := range req.Metrics {
		switch metric.Type {
		case "gauge":
			if metric.Gauge < 0 {
				return nil, status.Error(codes.InvalidArgument, "gauge value is negative")
			} else {
				s.db.SetGauge(metric.Name, metric.Gauge)
			}
		case "counter":
			if metric.Counter < 0 {
				return nil, status.Error(codes.InvalidArgument, "counter value is negative")
			} else {
				s.db.SetCounter(metric.Name, metric.Counter)
			}
		default:
			return nil, status.Error(codes.InvalidArgument, "invalid metric type")
		}
	}
	s.logger.Info("response to agent: success")
	return &response, nil
}
