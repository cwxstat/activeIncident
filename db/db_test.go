package db

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestNewActiveIncidentServer(t *testing.T) {
	tests := []struct {
		name    string
		want    *activeIncidentServer
		wantErr bool
	}{
		{
			name:    "",
			want:    &activeIncidentServer{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewActiveIncidentServer()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewActiveIncidentServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// err = got.addRecord()
			// if (err != nil) != tt.wantErr {
			// 	t.Errorf("addRecord() error = %v, wantErr %v", err, tt.wantErr)
			// 	return
			// }
		})
	}
}

func TestConn(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *mongo.Client
		wantErr bool
	}{
		{
			name:    "Simple connection test",
			args:    args{},
			want:    &mongo.Client{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Conn(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Conn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_ = got
		})
	}
}
