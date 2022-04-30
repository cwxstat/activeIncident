package wscrape

import (
	"encoding/json"
	"fmt"

	owm "github.com/briandowns/openweathermap"
	"github.com/cwxstat/activeIncident/dbutils"
)

var apiKey = dbutils.LookupEnv("OWM_API_KEY", "18ef17bf4ee75f4eafca0c158a33929b")

func Zips(zips []int) ([]string, error) {
	var out []string
	w, err := owm.NewCurrent("F", "EN", apiKey)
	if err != nil {
		return []string{}, err
	}
	sep := ""
	for _, zip := range zips {

		if err := w.CurrentByZip(zip, "US"); err != nil {
			return []string{}, err
		}
		b, err := json.Marshal(w)
		if err != nil {
			return []string{}, err
		}

		out = append(out, fmt.Sprintf("%s%s", sep, string(b)))
		sep = ","
	}
	return out, nil
}
