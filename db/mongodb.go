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

func (m *mongodb) entries(ctx context.Context) ([]activeIncidentEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	col := m.conn.Database("activeIncident").Collection("entries")
	cur, err := col.Find(ctx, bson.D{}, &options.FindOptions{
		Sort: map[string]interface{}{"_id": -1},
	})
	if err != nil {
		return nil, fmt.Errorf("mongodb.Find failed: %+v", err)
	}
	defer cur.Close(ctx)

	var out []activeIncidentEntry
	for cur.Next(ctx) {
		var v activeIncidentEntry
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

func (m *mongodb) addEntry(ctx context.Context, e activeIncidentEntry) error {
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
