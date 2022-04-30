package db

import (
	"context"
)

type database interface {
	entries(context.Context) ([]WeatherEntry, error)
	entriesMinutesAgo(context.Context, int) ([]WeatherEntry, error)
	addEntry(context.Context, WeatherEntry) error
	deleteAll(context.Context, string) error
	databaseCollection(string, string)
	disconnect(context.Context) error
}
