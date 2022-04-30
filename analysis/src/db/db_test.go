package db

import (
	"context"

	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestFull(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*30)
	defer cancel()
	as, err := NewActiveIncidentServer(ctx)
	if err != nil {
		t.FailNow()
	}

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
		TimeStamp:         db.NYtime(),
	})
	if err != nil {
		t.FailNow()
	}

	err = as.db.deleteAll(ctx, "Test Message")
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

func TestRetrieveEntries(t *testing.T) {
	as, err := NewActiveIncidentServer(context.TODO())
	if err != nil {
		t.FailNow()
	}

	result, err := as.db.entries(context.TODO(), 5)
	if err != nil {
		t.FailNow()
	}
	if len(result) < 5 {
		t.FailNow()
	}

}
