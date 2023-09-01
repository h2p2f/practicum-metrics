package grpcserver

import (
	"context"
	pb "github.com/h2p2f/practicum-metrics/proto"
	"go.uber.org/zap"
)

type Updater interface {
	SetGauge(name string, value float64)
	SetCounter(name string, value int64)
	GetCounter(name string) (value int64, err error)
	GetGauge(name string) (value float64, err error)
}

type Server struct {
	pb.UnimplementedMetricsServiceServer
	db     Updater
	logger *zap.Logger
}

func NewServer(db Updater, logger *zap.Logger) *Server {
	return &Server{db: db, logger: logger}
}

func (s *Server) UpdateMetric(
	ctx context.Context,
	req *pb.UpdateMetricRequest,
) (*pb.UpdateMetricResponse, error) {
	// Processing and validation of the received data
	var response pb.UpdateMetricResponse
	var err error
	s.logger.Info(
		"request from client:",
		zap.String("metric", req.Metric.Name),
		zap.String("type", req.Metric.Type),
		zap.Float64("gauge", req.Metric.Gauge),
		zap.Int64("counter", req.Metric.Counter))
	switch req.Metric.Type {
	case "gauge":
		if req.Metric.Gauge < 0 {
			response.Success = false
		} else {
			s.db.SetGauge(req.Metric.Name, req.Metric.Gauge)

			response.Metric = req.Metric
			response.Metric.Gauge, err = s.db.GetGauge(req.Metric.Name)
			if err != nil {
				response.Success = false
				response.Metric = nil
			}
			response.Success = true
		}
	case "counter":
		if req.Metric.Counter < 0 {
			response.Success = false
		} else {
			s.db.SetCounter(req.Metric.Name, req.Metric.Counter)
			response.Metric = req.Metric
			response.Metric.Counter, err = s.db.GetCounter(req.Metric.Name)
			if err != nil {
				response.Success = false
				response.Metric = nil
			}
			response.Success = true
		}
	default:
		response.Metric = nil
		response.Success = false
	}
	s.logger.Info("response from server:", zap.Bool("success", response.Success))
	return &response, err
}

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
				response.Success = false
			} else {
				s.db.SetGauge(metric.Name, metric.Gauge)
				response.Success = true
			}
		case "counter":
			if metric.Counter < 0 {
				response.Success = false
			} else {
				s.db.SetCounter(metric.Name, metric.Counter)
				response.Success = true
			}
		default:
			response.Success = false
		}
	}
	s.logger.Info("response to agent:", zap.Bool("success", response.Success))
	return &response, nil
}
