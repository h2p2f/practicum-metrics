package httpclient

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/h2p2f/practicum-metrics/internal/agent/compressor"
	"github.com/h2p2f/practicum-metrics/internal/agent/config"
	"github.com/h2p2f/practicum-metrics/internal/agent/models"
)

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

func SendJSONMetrics(logger *zap.Logger, data [][]byte, addr string) error {
	for _, d := range data {
		var metric models.Metric
		zipped, err := compressor.Compress(d)
		if err != nil {
			return err
		}
		client := resty.New()
		resp, err := client.R().
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

func SendBatchJSONMetrics(logger *zap.Logger, config *config.AgentConfig, data []byte, checkSum [32]byte) error {
	zipped, err := compressor.Compress(data)
	if err != nil {
		return err
	}
	//fmt.Println(fmt.Sprintf("%x", hash))
	hash := fmt.Sprintf("%x", checkSum)
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeaderVerbatim("HashSHA256", fmt.Sprintf("%x", hash)).
		SetBody(zipped).
		Post("http://" + config.ServerAddress + "/updates/")
	if err != nil {
		return err
	}
	logger.Info("response from server:",
		zap.Int("status code", resp.StatusCode()))
	return nil
}

func SendMetric(logger *zap.Logger, data []byte, checkSum [32]byte, address string) error {
	zipped, err := compressor.Compress(data)
	if err != nil {
		return err
	}
	client := resty.New()
	client.SetRetryCount(3).SetRetryWaitTime(1 * time.Second)
	hash := fmt.Sprintf("%x", checkSum)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeaderVerbatim("HashSHA256", hash).
		SetBody(zipped).
		Post("http://" + address + "/update/")
	if err != nil {
		return err
	}
	logger.Info("received response from server: ",
		zap.Int("status code", resp.StatusCode()))
	return nil
}
