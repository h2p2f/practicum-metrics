package httpclient

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/h2p2f/practicum-metrics/internal/agent/compressor"
	"github.com/h2p2f/practicum-metrics/internal/agent/models"
)

// Deprecated:SendJSONMetrics sends metrics to the server in JSON format one metric at a time.
// Required for backward compatibility. Currently not used.
func SendJSONMetrics(ctx context.Context, logger *zap.Logger, data [][]byte, addr string) error {
	for _, d := range data {
		var metric models.Metric
		zipped, err := compressor.Compress(d)
		if err != nil {
			return err
		}
		client := resty.New()
		resp, err := client.R().
			SetContext(ctx).
			SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetBody(zipped).
			SetResult(&metric).
			Post("http://" + addr + "/update/")
		if err != nil {
			return err
		}
		logger.Info("response from server:",
			zap.Int("status code", resp.StatusCode()),
			zap.String("metric", metric.ID))
	}
	return nil
}

// Deprecated:SendMetrics sends metrics to the server. Required for backward compatibility. Currently not used.
func SendMetrics(links []string) error {
	for _, link := range links {
		client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "text/plain").
			Post(link)
		if err != nil {
			return err
		}
		fmt.Print(resp)
	}
	return nil
}
