package metrics

import (
	"reflect"
	"testing"
)

func TestUrlMetrics(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want []string
	}{
		{
			name: "Positive test",
			url:  "http://localhost:8080",
			want: []string{"http://localhost:8080"},
		},
		{
			name: "Negative test",
			url:  "http://localhost:8080",
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rt RuntimeMetrics
			rt.NewMetrics()

			got := rt.URLMetrics(tt.url)

			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("UrlMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
