package db

import (
	"context"

	"fmt"
	"log"

	"os"
	"time"

	"github.com/cwxstat/activeIncident/constants"
	"github.com/cwxstat/activeIncident/dbutils"
	"github.com/cwxstat/activeIncident/weather/wscrape"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type weatherServer struct {
	db database
}

type WeatherEntry struct {
	Weather   []string  `json:"weather" bson:"weather"`
	TimeStamp time.Time `json:"date" bson:"date"`
}

func conn(ctx context.Context) (*mongo.Client, error) {

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Println("MONGO_URI environment variable not specified")
		return nil, fmt.Errorf("MONGO_URI environment variable not specified")
	}

	dbConn, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {

		log.Printf("failed to initialize connection to mongodb: %+v", err)
		return nil, err
	}
	if err := dbConn.Ping(ctx, readpref.Primary()); err != nil {
		log.Printf("ping to mongodb failed: %+v", err)
		return nil, err
	}

	return dbConn, nil

}

func NewWeather(ctx context.Context) (*weatherServer, error) {

	client, err := conn(ctx)
	if err != nil {
		return nil, err
	}

	w := &weatherServer{
		db: &mongodb{
			conn:       client,
			database:   dbutils.LookupEnv("MONGO_DB", "activeIncident"),
			collection: dbutils.LookupEnv("MONGO_WEATHER", "weather"),
		},
	}
	return w, nil
}

func (a *weatherServer) Disconnect(ctx context.Context) error {
	return a.db.disconnect(ctx)
}

func (a *weatherServer) AddEntry(ctx context.Context, entry *WeatherEntry) error {
	return a.db.addEntry(ctx, *entry)
}

func PopulateWeather() (*WeatherEntry, error) {
	we := &WeatherEntry{}
	var err error
	zips := []int{
		19027,
		18041,
		18426}
	we.Weather, err = wscrape.Zips(zips)
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
