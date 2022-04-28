package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongodb struct {
	conn *mongo.Client
}

func (m *mongodb) entries(ctx context.Context, minutes int) ([]ActiveIncidentEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	col := m.conn.Database("activeIncident").Collection("entries")

	// Only return these fields
	opts := options.Find().SetProjection(bson.D{
		{"incidents", 1},
		{"date", 1},
		{"_id", 1},
	})

	cur, err := col.Find(ctx,
		bson.D{{"date", bson.D{{"$gt", time.Now().Add(-time.Minute * time.Duration(minutes))}}}}, opts)
	if err != nil {
		return nil, fmt.Errorf("mongodb.Find failed: %+v", err)
	}
	defer cur.Close(ctx)

	var out []ActiveIncidentEntry
	for cur.Next(ctx) {
		var v ActiveIncidentEntry
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

func (m *mongodb) addEntry(ctx context.Context, e ActiveIncidentEntry) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	col := m.conn.Database("activeIncident").Collection("entries")
	if _, err := col.InsertOne(ctx, e); err != nil {
		return fmt.Errorf("mongodb.InsertOne failed: %+v", err)
	}
	return nil
}

func (m *mongodb) deleteAll(ctx context.Context, message string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	col := m.conn.Database("activeIncident").Collection("entries")
	if _, err := col.DeleteMany(ctx, bson.M{"message": message}); err != nil {
		return fmt.Errorf("mongodb.DeleteOne failed: %+v", err)
	}
	return nil
}
