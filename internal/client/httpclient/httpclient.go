package httpclient

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
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

type metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

// function Compress to compress data
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

func CompressMet(data []metrics) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	gz := gzip.NewWriter(buf)
	if err := json.NewEncoder(gz).Encode(data); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func SendMetrics(met [][]byte, address string) (err error) {
	for _, data := range met {
		//compress data, this comment wrote captain obvious
		buf, err := Compress(data)
		if err != nil {
			return err
		}
		//send data to server
		client := resty.New()
		//some autotests can be faster than server starts, so we need to retry three times, not more :)
		client.SetRetryCount(3).SetRetryWaitTime(200 * time.Millisecond)
		//upgrading request's headers
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetBody(buf).
			Post(address + "/update/")
		if err != nil {
			return err
		}
		logger.Log.Sugar().Info("received response from server: ", resp.StatusCode())
		fmt.Println("received response from server: ", resp.StatusCode())
	}
	return nil
}

func SendBatchMetrics(met []byte, address string) (err error) {
	//compress data, this comment wrote captain obvious
	//var dataToSend []byte
	//for _, data := range met {
	//	dataToSend = append(dataToSend, data...)
	//}
	//fmt.Println("dataToSend: ", dataToSend)
	//fmt.Println("met: ", met)
	buf, err := Compress(met)
	if err != nil {
		return err
	}
	//send data to server
	client := resty.New()
	//some autotests can be faster than server starts, so we need to retry three times, not more :)
	client.SetRetryCount(3).SetRetryWaitTime(200 * time.Millisecond)

	//upgrading request's headers
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(buf).
		Post(address + "/updates/")
	if err != nil {
		return err
	}
	logger.Log.Sugar().Info("received response from server: ", resp.StatusCode())
	return nil
}
