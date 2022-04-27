package scrape

import (
	"errors"
	"log"
	"time"

	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
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
			panic("Couldn't perform GET request to " + url + err.Error())
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

// Tag: returns station, incident, and error
func Tag(s string) ([]string, []string, error) {
	doc, err := html.Parse(strings.NewReader(s))
	station := []string{}
	incident := []string{}
	if err != nil {
		return station, incident, err
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {

					if strings.Contains(a.Val, "Lookup") {
						// fmt.Println(a.Val)
						station = append(station, a.Val)
					} else if strings.Contains(a.Val, "livecad") {
						incident = append(incident, a.Val)
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
	return station, incident, err
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
	s = strings.Replace(s, "LookupFD.asp?FDStation", "FDStation", -1)
	s = strings.Replace(s, "LookupEMS.asp?EMSStation", "EMSStation", -1)
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

func GetMainTable(s string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(s))
	r := []string{}
	stag := ""

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

					// Datetime has <br>
					if c.FirstChild.FirstChild.NextSibling != nil && c.FirstChild.FirstChild.NextSibling.Data == "br" {
						stag = c.FirstChild.FirstChild.Data + "T"
						stag = stag + c.FirstChild.FirstChild.NextSibling.NextSibling.Data
						r = append(r, stag)
					} else {
						if c.FirstChild.FirstChild.Data == "a" && c.FirstChild.FirstChild.FirstChild != nil {
							stag = c.FirstChild.FirstChild.FirstChild.Data
							r = append(r, stag)
						} else {
							stag = c.FirstChild.FirstChild.Data
							r = append(r, stag)
						}

					}

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

func GetTable(s string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(s))
	r := []string{}
	stime := ""

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

					// Datetime has <br>
					if c.FirstChild.FirstChild.NextSibling != nil && c.FirstChild.FirstChild.NextSibling.Data == "br" {
						stime = c.FirstChild.FirstChild.Data + "T"
						stime = stime + c.FirstChild.FirstChild.NextSibling.NextSibling.Data
						r = append(r, stime)
					} else {
						r = append(r, c.FirstChild.FirstChild.Data)
					}

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

type StationIncidentStatus struct {
	Time          time.Time
	Count         int64
	Station       map[string]string
	Incident      map[string]string
	Status        []string
	WebStatusPage string
	WebMainPage   string
}

type DB struct {
	Count       int64
	TimeUpdate  time.Time
	CreateTime  time.Time
	FileStation string
	Events      []StationIncidentStatus
}

func NewDB() *DB {
	return &DB{
		CreateTime: time.Now(),
		Events:     []StationIncidentStatus{},
	}
}

func (db *DB) GetsEverything() error {

	url := constants.WebCadURL + "livecad.asp?print=yes"
	mainPage, err := Get(url)
	if err != nil {
		return err
	}

	station, incident, err := Tag(mainPage)
	if err != nil {
		return err
	}

	for i, l := range incident {
		stationIncidentStatus := StationIncidentStatus{}
		statusPage, err := Get(GetDetail(l))
		if err != nil {
			return err
		}

		stationIncidentStatus.WebStatusPage = statusPage
		stationIncidentStatus.WebMainPage = mainPage
		if len(station) <= i {
			stationIncidentStatus.Time = time.Now()
			stationIncidentStatus.Station = map[string]string{"none": "none"}
			stationIncidentStatus.Incident = strip(l)

		} else {
			stationIncidentStatus.Time = time.Now()
			stationIncidentStatus.Station = strip(station[i])
			stationIncidentStatus.Incident = strip(l)
		}

		if status, err := GetTable(statusPage); err == nil {
			stationIncidentStatus.Status = status
		}
		db.WriteStationIncidentStatus("/data/activeIncidents.json", stationIncidentStatus)
		db.Events = append(db.Events, stationIncidentStatus)
		db.TimeUpdate = time.Now()
		db.Count++
	}
	return nil
}

func (db *DB) WriteStationIncidentStatus(file string, s StationIncidentStatus) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	f.WriteString(string(b))
	return nil
}

// WriteDB complete write
func (db *DB) WriteEventsAll(file string) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := json.Marshal(db.Events)
	if err != nil {
		return err
	}
	f.WriteString(string(b))
	return nil
}

func (db *DB) ClearDB(file string, count int64) error {
	if db.Count >= count {
		err := db.WriteEventsAll(file)
		if err != nil {
			log.Println(err)
			return err
		}
		db.Events = []StationIncidentStatus{}
		db.Count = 0
	}
	return nil
}
