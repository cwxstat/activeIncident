package db

import (
	"context"
	"fmt"
	"time"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongodb struct {
	conn       *mongo.Client
	database   string
	collection string
}

func (m *mongodb) databaseCollection(database string, collection string) {
	m.database = database
	m.collection = collection
}

func (m *mongodb) disconnect(ctx context.Context) error {
	return m.conn.Disconnect(ctx)
}

func (m *mongodb) entriesMinutesAgo(ctx context.Context, minutes int) ([]WeatherEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Only return these fields
	opts := options.Find().SetProjection(bson.D{
		{"incidents", 1},
		{"date", -1},
		{"_id", 1},
	})

	col := m.conn.Database(m.database).Collection(m.collection)

	cur, err := col.Find(ctx,
		bson.D{{"date", bson.D{{"$gt", time.Now().Add(-time.Minute * time.Duration(minutes))}}}}, opts)
	if err != nil {
		return nil, fmt.Errorf("mongodb.Find failed: %+v", err)
	}

	defer cur.Close(ctx)

	var out []WeatherEntry
	for cur.Next(ctx) {
		var v WeatherEntry
		if err := cur.Decode(&v); err != nil {
			return nil, fmt.Errorf("decoding mongodb record failed: %+v", err)
		}
		out = append(out, v)
	}
	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate on mongodb cursor: %+v", err)
	}
	return out, nil
}

func (m *mongodb) entries(ctx context.Context) ([]WeatherEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	col := m.conn.Database(m.database).Collection(m.collection)
	cur, err := col.Find(ctx, bson.D{}, &options.FindOptions{
		Sort: map[string]interface{}{"_id": -1},
	})
	if err != nil {
		return nil, fmt.Errorf("mongodb.Find failed: %+v", err)
	}
	defer cur.Close(ctx)

	var out []WeatherEntry
	for cur.Next(ctx) {
		var v WeatherEntry
		if err := cur.Decode(&v); err != nil {
			return nil, fmt.Errorf("decoding mongodb record failed: %+v", err)
		}
		out = append(out, v)
	}
	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate on mongodb cursor: %+v", err)
	}
	return out, nil
}

func (m *mongodb) addEntry(ctx context.Context, e WeatherEntry) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	col := m.conn.Database(m.database).Collection(m.collection)
	if _, err := col.InsertOne(ctx, e); err != nil {
		return fmt.Errorf("mongodb.InsertOne failed: %+v", err)
	}
	log.Printf("Added entry: %+v,\n ->%+v<-, ->%+v<-\n", e,m.database,m.collection)
	return nil
}

func (m *mongodb) deleteAll(ctx context.Context, message string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	col := m.conn.Database(m.database).Collection(m.collection)
	if _, err := col.DeleteMany(ctx, bson.M{"message": message}); err != nil {
		return fmt.Errorf("mongodb.DeleteOne failed: %+v", err)
	}

	return nil
}
