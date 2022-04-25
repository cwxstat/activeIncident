package db

import (
	"testing"
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
			got, err := NewActiveIncidentServer()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewActiveIncidentServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_ = got
		})
	}
}
