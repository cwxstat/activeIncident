package dbutils

import (
	"log"
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
