package db

import (
	"context"
)

type database interface {
	entries(context.Context, int) ([]ActiveIncidentEntry, error)
	addEntry(context.Context, ActiveIncidentEntry) error
	deleteAll(context.Context, string) error
}
