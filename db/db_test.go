package db

import (
	"context"

	"testing"
	"time"

	"github.com/cwxstat/activeIncident/db"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestFull(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancel()
	as, err := NewActiveIncidentServer(ctx)
	if err != nil {
		t.FailNow()
	}

	as.DatabaseCollection("test", "test")

	iwebp := []IncidentWebPage{
		IncidentWebPage{
			Page: "Page1",
		},
		IncidentWebPage{
			Page: "Page2",
		},
	}

	err = as.db.addEntry(ctx, ActiveIncidentEntry{
		MainWebPage:      "Main",
		IncidentWebPages: iwebp,
		Incidents:        []Incident{},
		Message:          "Test Message",
		TimeStamp:        db.NYtime(),
	})
	if err != nil {
		t.FailNow()
	}

	result, err := as.db.entriesMinutesAgo(ctx, 1)
	if err != nil || len(result) != 1 {
		t.FailNow()
	}

	err = as.db.deleteAll(ctx, "Test Message")
	if err != nil {
		t.FailNow()
	}

	err = as.db.disconnect(ctx)
	if err != nil {
		t.FailNow()
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
			name: "Simple connection test",
			args: args{
				ctx: context.TODO(),
			},
			want:    &mongo.Client{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(tt.args.ctx, time.Second*30)
		defer cancel()
		t.Run(tt.name, func(t *testing.T) {
			client, err := conn(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Conn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer client.Disconnect(ctx)
		})
	}
}

func TestNewActiveIncidentServer(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *activeIncidentServer
		wantErr bool
	}{
		{
			name: "NewActiveIncidentServer",
			args: args{
				ctx: context.TODO(),
			},
			want:    &activeIncidentServer{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewActiveIncidentServer(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewActiveIncidentServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_ = got
			//got.db.addEntry(context.Context, guestbookEntry)
		})
	}
}
