package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// activeIncidentEntry represents the message object returned in the API.
type activeIncidentEntry struct {
	Author  string    `json:"author" bson:"author"`
	Message string    `json:"message" bson:"message"`
	TimeStamp    time.Time `json:"date" bson:"date"`
}

type guestbookServer struct {
	db database
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

func AddRecord(client *mongo.Client) error {

	col := client.Database("guestbook").Collection("entries")

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()
	g := &activeIncidentEntry{
		Author:  "Susan and more",
		Message: "Here I ame",
		TimeStamp:    time.Now(),
	}
	v, err := col.InsertOne(ctx, g)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("no error: ", v)
	}
	return err
}

func (s *activeIncidentServer) addRecord() error {

	v := activeIncidentEntry{
		Author:  "Susan",
		Message: "Okay .. makes sense",
		TimeStamp:    time.Now(),
	}

	ctx := context.Background()
	connCtx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	if err := s.db.addEntry(connCtx, v); err != nil {
		return err
	}
	log.Printf("entry saved: author=%q message=%q", v.Author, v.Message)
	return nil
}

// main starts a server listening on $PORT responding to requests "GET
// /messages" and "POST /messages" with a JSON API.
func Mmain() {
	ctx := context.Background()

	// PORT environment variable is set in guestbook-backend.deployment.yaml.
	port := "8080"
	if port == "" {
		log.Fatal("PORT environment variable not specified")
	}

	mongoURI := os.Getenv("MONGO_URI")
	connCtx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	dbConn, err := mongo.Connect(connCtx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("failed to initialize connection to mongodb: %+v", err)
	}
	if err := dbConn.Ping(connCtx, readpref.Primary()); err != nil {
		log.Fatalf("ping to mongodb failed: %+v", err)
	}

	gs := &guestbookServer{
		db: &mongodb{
			conn: dbConn,
		},
	}

	log.Printf("backend server listening on port %s", port)
	http.Handle("/messages", gs)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func (s *guestbookServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("received request: method=%s path=%s", r.Method, r.URL.Path)
	if r.Method == http.MethodGet {
		s.getMessagesHandler(w, r)
	} else if r.Method == http.MethodPost {
		s.postMessageHandler(w, r)
	} else {
		http.Error(w, fmt.Sprintf("unsupported method %s", r.Method), http.StatusMethodNotAllowed)
	}
}

func (s *guestbookServer) getMessagesHandler(w http.ResponseWriter, r *http.Request) {
	entries, err := s.db.entries(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read entries: %+v", err), http.StatusInternalServerError)
		// TODO return JSON error
		return
	}
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		log.Printf("WARNING: failed to encode json into response: %+v", err)
	} else {
		log.Printf("%d entries returned", len(entries))
	}
}

func (s *guestbookServer) postMessageHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var v activeIncidentEntry
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request body into json: %+v", err), http.StatusBadRequest)
		return
	}
	if v.Author == "" {
		http.Error(w, "empty 'author' value", http.StatusBadRequest)
		return
	}
	if v.Message == "" {
		http.Error(w, "empty 'message' value", http.StatusBadRequest)
		return
	}

	v.TimeStamp = time.Now()

	if err := s.db.addEntry(r.Context(), v); err != nil {
		http.Error(w, fmt.Sprintf("failed to save entry: %+v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("entry saved: author=%q message=%q", v.Author, v.Message)
}
