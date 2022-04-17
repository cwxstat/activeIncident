package scrape

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"context"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cwxstat/activeIncident/constants"
	"golang.org/x/net/html"
)

var debug = false

// Headers contains all HTTP headers to send
var Headers = make(map[string]string)

// Cookies contains all HTTP cookies to send
var Cookies = make(map[string]string)

// SetDebug sets the debug status
// Setting this to true causes the panics to be thrown and logged onto the console.
// Setting this to false causes the errors to be saved in the Error field in the returned struct.
func SetDebug(d bool) {
	debug = d
}

// Header sets a new HTTP header
func Header(n string, v string) {
	Headers[n] = v
}

func Cookie(n string, v string) {
	Cookies[n] = v
}

// GetWithClient returns the HTML returned by the url using a provided HTTP client
func GetWithClient(url string, client *http.Client) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*800))
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		if debug {
			panic("Couldn't perform GET request to " + url)
		}
		return "", errors.New("couldn't perform GET request to " + url)
	}
	// Set headers
	for hName, hValue := range Headers {
		req.Header.Set(hName, hValue)
	}
	// Set cookies
	for cName, cValue := range Cookies {
		req.AddCookie(&http.Cookie{
			Name:  cName,
			Value: cValue,
		})
	}
	// Perform request
	resp, err := client.Do(req)
	if err != nil {
		if debug {
			panic("Couldn't perform GET request to " + url)
		}
		return "", errors.New("couldn't perform GET request to " + url)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if debug {
			panic("Unable to read the response body")
		}
		return "", errors.New("unable to read the response body")
	}
	return string(bytes), nil
}

type HTTP struct {
	client *http.Client
}

func Get(url string, client ...*http.Client) (string, error) {

	var newclient *http.Client
	if client == nil {
		newclient = &http.Client{}
	} else {
		newclient = client[0]
	}

	return GetWithClient(url, newclient)
}

type DB struct {
	r map[string]string
	v []string
}

func Tag(s string) ([]string, []string, error) {
	doc, err := html.Parse(strings.NewReader(s))
	r := []string{}
	l := []string{}
	if err != nil {
		return r, l, err
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {

					if strings.Contains(a.Val, "Lookup") {
						// fmt.Println(a.Val)
						r = append(r, a.Val)
					} else if strings.Contains(a.Val, "livecad") {
						l = append(l, a.Val)
					}

					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return r, l, err
}

func strip(s string) map[string]string {

	//fmt.Printf("%v\n", s)
	m := map[string]string{}
	s = cleanUp(s)
	for _, v := range strings.Split(s, "&") {
		ss := strings.Split(v, "=")
		if len(ss) == 2 {
			//fmt.Printf("M: %s, %s\n", ss[0], ss[1])
			m[ss[0]] = ss[1]
		}

	}
	return m
}

func cleanUp(s string) string {
	s = strings.Replace(s, "livecadcomments-fireems.asp?eid", "eid", -1)
	s = strings.Replace(s, "livecadcomments.asp?eid", "eid", -1)
	s = strings.Replace(s, "map.asp?type", "type", -1)
	s = strings.Replace(s, "<br>", " ", -1)
	s = strings.Replace(s, " @ ", " ", -1)
	return s
}

func GetDetail(purl string) string {
	url := constants.WebCadURL + purl
	return strings.Replace(url, " ", "%20", -1)
}

func GetTable(s string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(s))
	r := []string{}

	if err != nil {
		return r, err
	}
	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {

		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {

			if c.Data == "td" {

				if c.FirstChild.Data == "b" {
					//c = c.FirstChild
					return
				}

				if c.FirstChild.Data == "font" {
					r = append(r, c.FirstChild.FirstChild.Data)
				} else {
					r = append(r, c.FirstChild.Data)
				}

			}

			f(c)
		}
	}
	f(doc)

	return r, nil
}

func GetTableV2(s string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(s))
	r := []string{}

	if err != nil {
		return r, err
	}
	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {

		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {

			if c.Data == "td" {

				if c.FirstChild.Data == "b" {
					//c = c.FirstChild
					return
				}

				if c.FirstChild.Data == "font" {
					r = append(r, c.FirstChild.FirstChild.Data)
				} else {
					r = append(r, c.FirstChild.Data)
				}

			}

			f(c)
		}
	}
	f(doc)

	return r, nil
}

func BuildDb() ([]map[string]string, [][]string, error) {

	callTable := []map[string]string{}
	arriveTable := [][]string{}

	url := constants.WebCadURL + "livecad.asp?print=yes"
	r, err := Get(url)
	if err != nil {
		return nil, nil, err
	}

	result, link, err := Tag(r)
	if err != nil {
		return nil, nil, err
	}

	for _, result := range result {
		callTable = append(callTable, strip(result))
	}

	for _, l := range link {
		r, err = Get(GetDetail(l))
		if err != nil {
			return callTable, nil, err
		}

		arrive, err := GetTable(r)
		if err != nil {
			return callTable, nil, err
		}
		arriveTable = append(arriveTable, arrive)

	}

	return callTable, arriveTable, err

}

func Show() {
	c, a, err := BuildDb()
	if err != nil {
		log.Fatalf("No build")
	}
	for i, m := range c {
		for k, v := range m {
			fmt.Printf("%v: %v\n", k, v)
		}
		fmt.Printf("Status: %v\n\n", a[i])
	}
}

func ShowJson() {
	a, err := GetJson()
	if err != nil {
		log.Printf("Error in json")
	}
	println(string(a))
}

func WriteJson(filename string) error {

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		return err
	}
	defer f.Close()

	a, err := GetJson()
	if err != nil {
		log.Printf("Error in json")
	}
	if _, err = f.WriteString(string(a)); err != nil {
		return err
	}
	return nil
}

//TODO: Fix me
func GetJson() ([]byte, error) {
	c, a, err := BuildDb()
	if err != nil {
		log.Fatalf("No build")
	}

	return ToJson(c, a)

}

func ToJson(call []map[string]string, status [][]string) ([]byte, error) {
	type Calls struct {
		Call   map[string]string
		Status []string
	}

	calls := []*Calls{}

	if len(status) < len(call) {
		log.Printf("len(status) < len(call)\n")
		for i := len(status); i < len(call); i++ {
			status = append(status, []string{})
		}
	}

	for i, v := range call {
		nt := new(Calls)
		nt.Call = v
		nt.Status = status[i]
		calls = append(calls, nt)
	}

	type DB struct {
		Calls     []*Calls
		TimeStamp time.Time
	}

	return json.Marshal(DB{calls, time.Now()})

}
