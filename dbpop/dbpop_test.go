package dbpop

import (
	"testing"

	"github.com/cwxstat/activeIncident/db"
)

func TestAddDB(t *testing.T) {
	tests := []struct {
		name    string
		want    *db.ActiveIncidentEntry
		wantErr bool
	}{
		{
			name:    "Populate ActiveIncidentEntry",
			want:    &db.ActiveIncidentEntry{},
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
