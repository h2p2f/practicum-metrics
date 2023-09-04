// Package grpcclient implements the logic of sending metrics to the server.
package grpcclient

import (
	"context"
	pb "github.com/h2p2f/practicum-metrics/proto"
)

// GRPCSendMetric - send one metric to server, receive it from channel
func GRPCSendMetric(c pb.MetricsServiceClient, mCh <-chan *pb.Metric, done chan<- bool) error {

	var err error
	for m := range mCh {
		_, err := c.UpdateMetric(context.Background(), &pb.UpdateMetricRequest{Metric: m})
		if err != nil {
			return err
		}
		done <- true
	}
	return err
}

// GRPCSendMetrics - send batch metrics to server in one request
func GRPCSendMetrics(c pb.MetricsServiceClient, m []*pb.Metric) error {

	ctx := context.Background()
	_, err := c.UpdateMetrics(ctx, &pb.UpdateMetricsRequest{Metrics: m})

	return err
}
