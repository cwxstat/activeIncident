package cmd

import (
	"context"
	"testing"
	"time"

	"github.com/cwxstat/activeIncident/db"
	"github.com/cwxstat/activeIncident/dbinterface"
)

func TestPopulate(t *testing.T) {

	a, err := dbinterface.AddDB()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*30)
	defer cancel()
	as, err := db.NewActiveIncidentServer(ctx)
	if err != nil {
		t.FailNow()
	}

	err = as.AddEntry(ctx, a)
	if err != nil {
		t.FailNow()
	}
}
