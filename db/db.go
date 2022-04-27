package db

import (
	"context"

	"fmt"
	"log"

	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	DispatchTime    string `json:"dispatchTime" bson:"dispatchTime"`
	Station         string `json:"station" bson:"station"`
	IncidentStatus  []IncidentStatus
}

// activeIncidentEntry represents the message object returned in the API.
type activeIncidentEntry struct {
	MainWebPage      string `json:"mainWebPage" bson:"mainWebPage"`
	IncidentWebPages []IncidentWebPage
	Incidents        []Incident
	Message          string    `json:"message" bson:"message"`
	TimeStamp        time.Time `json:"date" bson:"date"`
}

type activeIncidentServer struct {
	db database
}

/*
    connCtx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
*/
func conn(ctx context.Context) (*mongo.Client, error) {

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Println("MONGO_URI environment variable not specified")
		return nil, fmt.Errorf("MONGO_URI environment variable not specified")
	}

	dbConn, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {

		log.Printf("failed to initialize connection to mongodb: %+v", err)
		return nil, err
	}
	if err := dbConn.Ping(ctx, readpref.Primary()); err != nil {
		log.Printf("ping to mongodb failed: %+v", err)
		return nil, err
	}

	return dbConn, nil

}

func NewActiveIncidentServer(ctx context.Context) (*activeIncidentServer, error) {

	client, err := conn(ctx)
	if err != nil {
		return nil, err
	}
	a := &activeIncidentServer{
		db: &mongodb{
			conn: client,
		},
	}
	return a, nil
}

func (s *activeIncidentServer) addRecord() error {

	v := activeIncidentEntry{
		MainWebPage: "Susan",
		Message:     "Okay .. makes sense",
		TimeStamp:   time.Now(),
	}

	ctx := context.Background()
	connCtx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	if err := s.db.addEntry(connCtx, v); err != nil {
		return err
	}
	log.Printf("entry saved: author=%q message=%q", v.MainWebPage, v.Message)
	return nil
}
