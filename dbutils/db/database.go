package db

import (
	"context"
)

type Database interface {
	Entries(context.Context) ([]interface{}, error)
	EntriesMinutesAgo(context.Context, int) (interface{}, error)
	AddEntry(context.Context, interface{}) error
	DeleteAll(context.Context, string) error
	DatabaseCollection(string, string)
	Disconnect(context.Context) error
}
