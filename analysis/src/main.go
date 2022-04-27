package main

import (
	"context"
	"example.com/m/v2/db"
)

func main() {

	as, err := db.NewActiveIncidentServer(context.TODO())
	if err != nil {
		panic(err)
	}
	as.GetEntries(context.TODO(), "2022-04-27T14:46:51.143+00:00")

}
