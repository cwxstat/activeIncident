/*
https://webapp07.montcopa.org/eoc/cadinfo/livecad.asp?print=yes

*/

package scrape

import (
	"fmt"
	"log"
	"strings"
	"testing"

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

	for i, _ := range station {
		fmt.Printf("%v\n", strip(station[i]))
		fmt.Printf("%v\n", strip(incident[i]))
	}

	for i, l := range incident {
		util.WriteString(fmt.Sprintf("GetDetail%d", i), r, 0644)
		r, err = Get(GetDetail(l))
		if err != nil {
			t.Fatalf("err: %s\n", err)
		}

		if status, err := GetTable(r); err == nil {
			fmt.Printf("%v\n", status)
		}

	}

}

func TestGetBuildDB(t *testing.T) {

	c, a, err := BuildDb()
	if err != nil {
		t.FailNow()
	}
	for i, m := range c {
		for k, v := range m {
			fmt.Printf("%v: %v\n", k, v)
		}
		fmt.Printf("Status: %v\n\n", a[i])
	}

}

func TestShow(t *testing.T) {
	Show()
}

func TestGetTableV2(t *testing.T) {

	a, _ := GetTableV2(test_fixtures.Detail())
	for _, v := range a {
		fmt.Printf("%v\n", v)
	}
}

func Test_GetJson(t *testing.T) {
	a, err := GetJson()
	if err != nil {
		t.FailNow()
	}
	println(string(a))
}

func Test_WriteJson(t *testing.T) {
	WriteJson("./testfile")

}
