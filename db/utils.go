package db

import (
	"log"
	"time"
)

func NYtime() string {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Printf("Error loading location: %v", err)
		return time.Now().Format("2006-01-02T15:04:05")
	}
	return time.Now().In(loc).Format("2006-01-02T15:04:05")

}
