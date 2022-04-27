package scrape

import (
	"errors"

	"context"
	"io/ioutil"
	"net/http"
	"time"

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
