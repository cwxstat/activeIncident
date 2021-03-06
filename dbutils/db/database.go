package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database interface {
	Entries(context.Context) ([]interface{}, error)
	EntriesMinutesAgo(context.Context, int, *options.FindOptions) (interface{}, error)
	AddEntry(context.Context, interface{}) error
	DeleteAll(context.Context, string) error
	DatabaseCollection(string, string)
	Disconnect(context.Context) error
}
