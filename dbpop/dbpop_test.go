package dbpop

import (
	"testing"

	"github.com/cwxstat/activeIncident/active"
)

func TestAddDB(t *testing.T) {
	tests := []struct {
		name    string
		want    *active.ActiveIncidentEntry
		wantErr bool
	}{
		{
			name:    "Populate ActiveIncidentEntry",
			want:    &active.ActiveIncidentEntry{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PopulateActiveIncidentEntry()
			if (err != nil) != tt.wantErr {
				t.Errorf("AddDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_ = got
		})
	}
}
