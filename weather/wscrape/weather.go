package wscrape

import (
	"encoding/json"

	"time"

	owm "github.com/briandowns/openweathermap"
	"github.com/cwxstat/activeIncident/dbutils"
)

var apiKey = dbutils.LookupEnv("OWM_API_KEY", "18ef17bf4ee75f4eafca0c158a33929b")

type Weather struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Message int    `json:"message"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Base    string `json:"base"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Main struct {
		Temp      float64 `json:"temp"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		FeelsLike float64 `json:"feels_like"`
		Pressure  int     `json:"pressure"`
		SeaLevel  int     `json:"sea_level"`
		GrndLevel int     `json:"grnd_level"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Rain struct {
		ThreeH int `json:"3h"`
	} `json:"rain"`
	Snow struct {
		ThreeH int `json:"3h"`
	} `json:"snow"`
	Dt       int    `json:"dt"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
	Timezone int    `json:"timezone"`
	Unit     string `json:"Unit"`
	Lang     string `json:"Lang"`
	Key      string `json:"Key"`
}

type WeatherResponse struct {
	Date    time.Time `json:"date"`
	Weather []Weather `json:"weather"`
}

func Zips(zips []int) (WeatherResponse, error) {
	var WeatherResponse = WeatherResponse{}
	WeatherResponse.Date = dbutils.NYtime()
	w, err := owm.NewCurrent("F", "EN", apiKey)
	if err != nil {
		return WeatherResponse, err
	}

	for _, zip := range zips {

		if err := w.CurrentByZip(zip, "US"); err != nil {
			return WeatherResponse, err
		}

		b, err := json.Marshal(w)
		if err != nil {
			return WeatherResponse, err
		}
		weather := &Weather{}
		json.Unmarshal(b, weather)

		WeatherResponse.Weather = append(WeatherResponse.Weather, *weather)

	}
	return WeatherResponse, nil
}
