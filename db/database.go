package db

import (
	"context"
)

type database interface {
	entries(context.Context) ([]ActiveIncidentEntry, error)
	addEntry(context.Context, ActiveIncidentEntry) error
	deleteAll(context.Context, string) error
	databaseCollection(string, string)
	disconnect(context.Context) error
}
