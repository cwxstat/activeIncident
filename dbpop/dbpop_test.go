package dbpop

import (
	"bytes"
	"encoding/gob"
	"github.com/cwxstat/activeIncident/fixtures"
	"reflect"
	"testing"

	"github.com/cwxstat/activeIncident/active"
)

func TestAddDB(t *testing.T) {
	tests := []struct {
		name    string
		want    *active.ActiveIncidentEntry
		wantErr bool
	}{
		{
			name:    "Populate ActiveIncidentEntry",
			want:    &active.ActiveIncidentEntry{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PopulateActiveIncidentEntry()
			if (err != nil) != tt.wantErr {
				t.Errorf("AddDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_ = got
		})
	}
}

func Test_populateIncident(t *testing.T) {

	var network bytes.Buffer // Stand-in for a network connection
	dec := gob.NewDecoder(&network)
	path := fixtures.Path("./testdata/scrapeMainTable.enc")
	buf, err := fixtures.Read(path)
	n, err := network.Write(buf)
	if err != nil || n == 0 {
		t.Fatal(err)
	}

	var list []string
	err = dec.Decode(&list)
	if err != nil {
		t.Fatal(err)
	}

	result, err := populateIncident(list)
	if err != nil {
		t.Fatal(err)
	}
	e, err := ExpectedValue()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, e) {
		t.Fail()
	}
}

func ExpectedValue() ([]active.Incident, error) {
	var network bytes.Buffer // Stand-in for a network connection
	a := []active.Incident{}
	dec := gob.NewDecoder(&network)
	path := fixtures.Path("./testdata/expectedResult.enc")
	buf, err := fixtures.Read(path)
	n, err := network.Write(buf)
	if err != nil || n == 0 {
		return a, err
	}
	err = dec.Decode(&a)
	return a, err
}
