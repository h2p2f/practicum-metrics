package database

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/h2p2f/practicum-metrics/internal/server/mocks"
	"testing"
)

// this is a meaningless test since functions either do not return
// a result or have no input parameters other than the context
func TestWrite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockDataBaser(ctrl)
	m.EXPECT().
		Write(gomock.Any(), gomock.Any()).
		Return(nil)

	{
		tests := []struct {
			name    string
			wString string
			wantErr bool
		}{
			{
				name:    "TestWrite",
				wString: "{\"id\":\"OtherSys\",\"type\":\"gauge\",\"value\":1070588}",
				wantErr: false,
			},
		}
		for _, tt := range tests {
			var arg [][]byte
			ctx := context.TODO()
			arg = append(arg, []byte(tt.wString))
			err := m.Write(ctx, arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Error: %v", err)
			}
		}
	}
}

func TestRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockDataBaser(ctrl)
	m.EXPECT().
		Read(gomock.Any()).
		Return(nil, nil)

	{
		tests := []struct {
			name    string
			want    [][]byte
			wantErr bool
		}{
			{
				name:    "TestRead",
				want:    nil,
				wantErr: false,
			},
		}
		for _, tt := range tests {
			ctx := context.TODO()
			got, err := m.Read(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Error: %v", err)
			}
			if got != nil {
				t.Errorf("Got: %v", got)
			}
		}
	}
}
