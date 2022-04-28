package main

import (
   "context"
   "fmt"
   "os"
   "sync"
   "log"

   "go.mongodb.org/mongo-driver/bson"
   "go.mongodb.org/mongo-driver/mongo"
   "go.mongodb.org/mongo-driver/mongo/options"
   "go.mongodb.org/mongo-driver/mongo/readpref"
)

func iterateChangeStream(routineCtx context.Context, waitGroup sync.WaitGroup, stream *mongo.ChangeStream) {
   defer stream.Close(routineCtx)
   defer waitGroup.Done()
   for stream.Next(routineCtx) {
      var data bson.M
      if err := stream.Decode(&data); err != nil {
            panic(err)
      }
      fmt.Printf("%v\n", data)
   }
}

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

func main() {
   client, err := conn(context.TODO())
   if err != nil {
      panic(err)
   }
   defer client.Disconnect(context.TODO())

   database := client.Database("activeIncident")
   episodesCollection := database.Collection("entries")

   var waitGroup sync.WaitGroup

   episodesStream, err := episodesCollection.Watch(context.TODO(), mongo.Pipeline{})
   if err != nil {
      panic(err)
   }
   waitGroup.Add(1)
   routineCtx, cancelFn := context.WithCancel(context.Background())
   defer cancelFn()
   go iterateChangeStream(routineCtx, waitGroup, episodesStream)

   waitGroup.Wait()
}
