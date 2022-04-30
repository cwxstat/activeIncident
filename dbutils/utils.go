package dbutils

import (
	"log"
	"os"
	"time"
)

func NYtime() time.Time {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Printf("Error loading location: %v", err)
		return time.Now()
	}
	return time.Now().In(loc)

}

func LookupEnv(key string, defaultValue string) string {
	env := defaultValue
	if val, ok := os.LookupEnv(key); ok {
		env = val
	}
	return env
}
