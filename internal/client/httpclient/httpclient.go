package httpclient

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/h2p2f/practicum-metrics/internal/logger"
	"time"
)

func init() {
	if err := logger.InitLogger("info"); err != nil {
		fmt.Println(err)
	}
}

// Compress function to compress data
func Compress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	gz := gzip.NewWriter(buf)
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// SendMetric  - this function sends one metric to server
func SendMetric(met []byte, checkSum [32]byte, address string) (err error) {

	buf, err := Compress(met)
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
		SetBody(buf).
		Post("http://" + address + "/update/")
	if err != nil {
		return err
	}
	logger.Log.Sugar().Info("received response from server: ", resp.StatusCode())
	return nil
}

// SendBatchMetrics - this function sends all metrics to server
func SendBatchMetrics(met []byte, checkSum [32]byte, address string) (err error) {
	//compress data
	buf, err := Compress(met)
	if err != nil {
		return err
	}
	//send data to server
	client := resty.New()
	//some autotests can be faster than server starts, so we need to retry three times, not more :)
	client.SetRetryCount(3).SetRetryWaitTime(200 * time.Millisecond)
	hash := fmt.Sprintf("%x", checkSum)
	//upgrading request's headers
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeaderVerbatim("HashSHA256", hash).
		SetBody(buf).
		Post("http://" + address + "/updates/")
	if err != nil {
		return err
	}
	logger.Log.Sugar().Info("received response from server: ", resp.StatusCode())
	return nil
}
