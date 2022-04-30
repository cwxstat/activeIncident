package dbutils

import (
	"log"
	"os"
	"strings"
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
		env = strings.Replace(val, "\n", "", -1)
	}
	return env
}
