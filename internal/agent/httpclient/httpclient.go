// Package httpclient реализует логику отправки метрик на сервер.
//
// Package httpclient implements the logic of sending metrics to the server.
package httpclient

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/h2p2f/practicum-metrics/internal/agent/compressor"
	"github.com/h2p2f/practicum-metrics/internal/agent/config"
	"github.com/h2p2f/practicum-metrics/internal/agent/models"
)

// SendMetrics отправляет метрики на сервер. Необходима для обратной совместимости. В данный момент не используется.
//
// SendMetrics sends metrics to the server. Required for backward compatibility. Currently not used.
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

// SendJSONMetrics отправляет метрики на сервер в формате JSON по одной метрике за раз.
// Необходима для обратной совместимости. В данный момент не используется.
//
// SendJSONMetrics sends metrics to the server in JSON format one metric at a time.
// Required for backward compatibility. Currently not used.
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

// SendBatchJSONMetrics отправляет метрики на сервер в формате JSON в пакетном режиме.
//
// SendBatchJSONMetrics sends metrics to the server in JSON format in batch mode.
func SendBatchJSONMetrics(logger *zap.Logger, config *config.AgentConfig, data []byte, checkSum [32]byte) error {
	toSend, err := compressor.Compress(data)
	if err != nil {
		return err
	}
	hash := fmt.Sprintf("%x", checkSum)
	if config.PublicKey != nil {
		toSend, err = rsa.EncryptPKCS1v15(rand.Reader, config.PublicKey, toSend)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Encoding", "gzip").
		SetHeaderVerbatim("HashSHA256", fmt.Sprintf("%x", hash)).
		SetBody(toSend).
		Post("http://" + config.ServerAddress + "/updates/")
	if err != nil {
		return err
	}
	logger.Info("response from server:",
		zap.Int("status code", resp.StatusCode()))
	return nil
}

// SendMetric отправляет метрику на сервер. Используется воркерами.
//
// SendMetric sends a metric to the server. Used by workers.
func SendMetric(logger *zap.Logger, data []byte, checkSum [32]byte, config *config.AgentConfig) error {
	toSend, err := compressor.Compress(data)
	if err != nil {
		return err
	}
	client := resty.New()
	client.SetRetryCount(3).SetRetryWaitTime(1 * time.Second)
	hash := fmt.Sprintf("%x", checkSum)
	if config.PublicKey != nil {
		toSend, err = rsa.EncryptPKCS1v15(rand.Reader, config.PublicKey, toSend)
		if err != nil {
			return err
		}
	}
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeaderVerbatim("HashSHA256", hash).
		SetBody(toSend).
		Post("http://" + config.ServerAddress + "/update/")
	if err != nil {
		return err
	}
	logger.Info("received response from server: ",
		zap.Int("status code", resp.StatusCode()))
	return nil
}
