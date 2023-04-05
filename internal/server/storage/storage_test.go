package storage

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestSetGaugeStorage(t *testing.T) {

	tests := []struct {
		name   string
		metric string
		value  float64
		want   *MemStorage
	}{
		{
			name:   "Positive test",
			metric: "CPU",
			value:  0.0001,
			want: &MemStorage{
				Gauges: map[string][]float64{
					"CPU": {0.0001},
				},
				Counters: nil,
			},
		},
		{
			name:   "Negative test",
			metric: "CPU",
			value:  78,
			want: &MemStorage{
				Gauges: map[string][]float64{
					"CPU": {0.0001},
				},
				Counters: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMemStorage()
			got.SetGauge(tt.metric, tt.value)
			if got == tt.want {
				t.Errorf("NewStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetCounterStorage(t *testing.T) {

	tests := []struct {
		name   string
		metric string
		value  int64
		want   *MemStorage
	}{
		{
			name:   "Positive test",
			metric: "CPU",
			value:  200,
			want: &MemStorage{
				Gauges: nil,
				Counters: map[string]int64{
					"CPU": 200,
				},
			},
		},
		{
			name:   "Negative test",
			metric: "CPU",
			value:  78,
			want: &MemStorage{
				Gauges: nil,
				Counters: map[string]int64{
					"CPU": 200,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMemStorage()
			got.SetCounter(tt.metric, tt.value)
			if got == tt.want {
				t.Errorf("NewStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGaugeStorage(t *testing.T) {

	tests := []struct {
		name   string
		metric string
		want   string
	}{
		{
			name:   "Positive test",
			metric: "CPU",
			want:   "0.0001",
		},
		{
			name:   "Negative test",
			metric: "CPU",
			want:   "0.0001",
		},
	}
	got := NewMemStorage()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got.SetGauge(tt.metric, 0.0001)
			val, _ := got.GetGauge(tt.metric)
			res, _ := strconv.ParseFloat(tt.want, 64)
			if val[0] != res {
				t.Errorf("NewStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCounterStorage(t *testing.T) {

	tests := []struct {
		name   string
		metric string
		want   int64
	}{
		{
			name:   "Positive test 1",
			metric: "CPU",
			want:   200,
		},
		{
			name:   "Positive test 2",
			metric: "Memory",
			want:   200,
		},
	}
	got := NewMemStorage()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got.SetCounter(tt.metric, 200)
			val, _ := got.GetCounter(tt.metric)
			if val != tt.want {
				t.Errorf("NewStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllGaugesStorage(t *testing.T) {

	tests := []struct {
		name string
		want map[string][]float64
	}{
		{
			name: "Positive test 1",
			want: map[string][]float64{
				"CPU": {0.0001},
			},
		},
		{
			name: "Positive test 2",
			want: map[string][]float64{
				"CPU": {0.0001, 0.0001},
			},
		},
	}
	got := NewMemStorage()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got.SetGauge("CPU", 0.0001)
			val := got.GetAllGauges()
			assert.Equal(t, val, tt.want)
		})
	}
}
