package db

import (
	"context"
	"fmt"

	"time"

	"github.com/cwxstat/activeIncident/dbutils"
	db2 "github.com/cwxstat/activeIncident/dbutils/db"
)

type IncidentWebPage struct {
	Page string `json:"incidentPage" bson:"incidentPage"`
}

type IncidentStatus struct {
	TimeStamp string `json:"timeStamp" bson:"timeStamp"`
	Unit      string `json:"unit" bson:"unit"`
	Status    string `json:"status" bson:"status"`
	Notes     string `json:"notes" bson:"notes"`
}

type Incident struct {
	IncidentNo      string `json:"incidentNo" bson:"incidentNo"`
	IncidentType    string `json:"incidentType" bson:"incidentType"`
	IncidentSubTupe string `json:"incidentSubType" bson:"incidentSubType"`
	Location        string `json:"location" bson:"location"`
	Municipality    string `json:"municipality" bson:"municipality"`
	DispatchTime    string `json:"dispatchTime" bson:"dispatchTime"`
	Station         string `json:"station" bson:"station"`
	IncidentStatus  []IncidentStatus
}

// ActiveIncidentEntry represents the message object returned in the API.
type ActiveIncidentEntry struct {
	MainWebPage      string `json:"mainWebPage" bson:"mainWebPage"`
	IncidentWebPages []IncidentWebPage
	Incidents        []Incident
	Message          string    `json:"message" bson:"message"`
	TimeStamp        time.Time `json:"date" bson:"date"`
}

type activeIncidentServer struct {
	db db2.Database
}

func NewActiveIncidentServer(ctx context.Context) (*activeIncidentServer, error) {

	client, err := dbutils.Conn(ctx)
	if err != nil {
		return nil, err
	}

	a := &activeIncidentServer{
		db: &db2.Mongodb{
			Conn:       client,
			Database:   dbutils.LookupEnv("MONGO_DB", "activeIncident"),
			Collection: dbutils.LookupEnv("MONGO_COLLECTION", "events"),
		},
	}
	return a, nil
}

func (a *activeIncidentServer) Disconnect(ctx context.Context) error {
	return a.db.Disconnect(ctx)
}

func (a *activeIncidentServer) AddEntry(ctx context.Context, entry *ActiveIncidentEntry) error {
	return a.db.AddEntry(ctx, *entry)
}

func (a *activeIncidentServer) EntriesMinutesAgo(ctx context.Context, min int) ([]ActiveIncidentEntry, error) {
	result, err := a.db.EntriesMinutesAgo(ctx, min)
	if err != nil {
		return nil, err
	}

	if val, ok := result.([]ActiveIncidentEntry); ok {
		return val, nil
	}
	return nil, fmt.Errorf("EntriesMinutesAgo")
}

func (a *activeIncidentServer) DatabaseCollection(database string, collection string) {
	a.db.DatabaseCollection(database, collection)
}
