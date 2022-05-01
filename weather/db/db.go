package db

import (
	"context"

	"log"

	"time"

	"github.com/cwxstat/activeIncident/constants"
	"github.com/cwxstat/activeIncident/dbutils"
	"github.com/cwxstat/activeIncident/dbutils/db"
	"github.com/cwxstat/activeIncident/weather/wscrape"
)

type weatherServer struct {
	db db.Database
}

type WeatherEntry struct {
	WeatherResponse wscrape.WeatherResponse
	TimeStamp       time.Time `json:"date" bson:"date"`
}

func NewWeather(ctx context.Context) (*weatherServer, error) {

	client, err := dbutils.Conn(ctx)
	if err != nil {
		return nil, err
	}

	w := &weatherServer{
		db: &db.Mongodb{
			Conn:       client,
			Database:   dbutils.LookupEnv("MONGO_DB", "activeIncident"),
			Collection: dbutils.LookupEnv("MONGO_WEATHER", "weather"),
		},
	}
	return w, nil
}

func (a *weatherServer) Disconnect(ctx context.Context) error {
	return a.db.Disconnect(ctx)
}

func (a *weatherServer) AddEntry(ctx context.Context, entry *WeatherEntry) error {
	return a.db.AddEntry(ctx, *entry)
}

func PopulateWeather() (*WeatherEntry, error) {
	we := &WeatherEntry{}
	var err error
	zips := []int{
		19027,
		18041,
		18426}
	we.WeatherResponse, err = wscrape.Zips(zips)
	we.TimeStamp = dbutils.NYtime()
	if err != nil {
		return &WeatherEntry{}, err
	}
	return we, nil
}

// go RunInGoRoutine()
func RunInGoRoutine(countlimit ...int64) error {

	for {

		err := func() error {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*40)
			defer cancel()
			ais, err := NewWeather(ctx)
			if err != nil {
				log.Println(err)
				time.Sleep(constants.ErrorBackoff)
				return err
			}

			a, err := PopulateWeather()
			if err != nil {
				log.Println(err)
				return err
			}

			err = ais.AddEntry(ctx, a)
			if err != nil {
				log.Println(err)
				return err
			}
			if err := ais.Disconnect(ctx); err != nil {
				log.Println("as.Disconnect: ", err)
			}
			log.Println("Weather data added")
			return nil

		}()
		// exit for testing
		if len(countlimit) > 0 {
			return err
		}
		time.Sleep(constants.RefreshRate)
	}

}
