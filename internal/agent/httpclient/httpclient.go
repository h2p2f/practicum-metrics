// Package httpclient implements the logic of sending metrics to the server.
package httpclient

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/h2p2f/practicum-metrics/internal/agent/compressor"
	"github.com/h2p2f/practicum-metrics/internal/agent/config"
)

// SendBatchJSONMetrics sends metrics to the server in JSON format in batch mode.
func SendBatchJSONMetrics(ctx context.Context, logger *zap.Logger, config *config.AgentConfig, data []byte, checkSum [32]byte) error {
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

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	client := resty.New()

	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("X-Real-IP", config.IPaddr.String()).
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

// SendMetric sends a metric to the server. Used by workers.
func SendMetric(
	ctx context.Context,
	logger *zap.Logger,
	data []byte,
	checkSum [32]byte,
	config *config.AgentConfig) error {

	toSend, err := compressor.Compress(data)
	if err != nil {
		return err
	}
	client := resty.New()
	client.SetRetryCount(config.RetryCount).SetRetryWaitTime(config.RetryWaitTime * time.Second)
	hash := fmt.Sprintf("%x", checkSum)
	if config.PublicKey != nil {
		toSend, err = rsa.EncryptPKCS1v15(rand.Reader, config.PublicKey, toSend)
		if err != nil {
			return err
		}
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("X-Real-IP", config.IPaddr.String()).
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
