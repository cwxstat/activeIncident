package cmd

import (
	"context"
	"testing"
	"time"

	"github.com/cwxstat/activeIncident/active"
	"github.com/cwxstat/activeIncident/dbpop"
)

func TestPopulate(t *testing.T) {

	a, err := dbpop.PopulateActiveIncidentEntry()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*30)
	defer cancel()
	as, err := active.NewActiveIncidentServer(ctx)
	if err != nil {
		t.FailNow()
	}

	err = as.AddEntry(ctx, a)
	if err != nil {
		t.FailNow()
	}
}
