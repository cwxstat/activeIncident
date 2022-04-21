/*
https://webapp07.montcopa.org/eoc/cadinfo/livecad.asp?print=yes

*/

package scrape

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/cwxstat/activeIncident/constants"
	test_fixtures "github.com/cwxstat/activeIncident/test-fixtures"
	"github.com/mchirico/tlib/util"
	"golang.org/x/net/html"
)

func TestBegin(t *testing.T) {
	s := test_fixtures.Page()
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					fmt.Println(a.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
}

func TestTag(t *testing.T) {
	defer util.NewTlib().ConstructDir()()

	result, _, err := Tag(test_fixtures.Page())
	if err != nil {
		t.FailNow()
	}

	strip(result[0])

}

func TestGetTable(t *testing.T) {
	r := test_fixtures.Table()
	result, err := GetTable(r)
	if err != nil {
		t.FailNow()
	}
	fmt.Printf("%v\n", result)
}

func Test_LiveCheck(t *testing.T) {

	defer util.NewTlib().ConstructDir()()

	url := constants.WebCadURL + "livecad.asp?print=yes"
	r, err := Get(url)
	if err != nil {
		t.Fatalf("err: %s\n", err)
	}

	util.WriteString("mainPage", r, 0644)

	station, incident, err := Tag(r)
	if err != nil {
		t.FailNow()
	}

	util.WriteString("mainPage", r, 0644)

	for i, l := range incident {
		util.WriteString(fmt.Sprintf("GetDetail%d", i), r, 0644)
		r, err = Get(GetDetail(l))
		if err != nil {
			t.Fatalf("err: %s\n", err)
		}

		if len(station) <= i {
			fmt.Printf("station: %v, incident: %v\n", "none", strip(l))
		} else {

			fmt.Printf("station: %v, incident: %v\n", strip(station[i]), strip(l))
		}

		if status, err := GetTable(r); err == nil {
			fmt.Printf("%v\n", status)
		}

	}

}

func TestDB_GetsEverything(t *testing.T) {
	type fields struct {
		Time   time.Time
		Events []StationIncidentStatus
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "One run",
			fields:  fields{
				Time:   time.Time{},
				Events: []StationIncidentStatus{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewDB()
			if err := db.GetsEverything(); (err != nil) != tt.wantErr {
				t.Errorf("DB.GetsEverything() error = %v, wantErr %v", err, tt.wantErr)
			}
			db.WriteDB()
		})
	}
}
