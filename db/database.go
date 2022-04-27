package db

import (
	"context"
)

type database interface {
	entries(context.Context) ([]activeIncidentEntry, error)
	addEntry(context.Context, activeIncidentEntry) error
	deleteAll(context.Context, string) error
}
