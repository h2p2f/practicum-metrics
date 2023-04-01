package storage

import (
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
