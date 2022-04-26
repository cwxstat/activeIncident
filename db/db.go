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

type IncidentPage struct {
	IncidentPage string `json:"incidentPage" bson:"incidentPage"`
}

// activeIncidentEntry represents the message object returned in the API.
type activeIncidentEntry struct {
	MainPage      string `json:"mainPage" bson:"mainPage"`
	IncidentPages []IncidentPage
	Message       string    `json:"message" bson:"message"`
	TimeStamp     time.Time `json:"date" bson:"date"`
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
		MainPage:  "Susan",
		Message:   "Okay .. makes sense",
		TimeStamp: time.Now(),
	}

	ctx := context.Background()
	connCtx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	if err := s.db.addEntry(connCtx, v); err != nil {
		return err
	}
	log.Printf("entry saved: author=%q message=%q", v.MainPage, v.Message)
	return nil
}
