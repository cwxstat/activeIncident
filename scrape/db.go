package scrape

import (
	"time"

	"github.com/cwxstat/activeIncident/constants"
	"github.com/cwxstat/activeIncident/db"
)

func AddDB() (*db.ActiveIncidentEntry, error) {

	aie := &db.ActiveIncidentEntry{}
	url := constants.WebCadURL + "livecad.asp?print=yes"
	r, err := Get(url)
	if err != nil {
		return aie, err
	}
	aie.MainWebPage = r
	aie.TimeStamp = time.Now()
	aie.Incidents, err = PopulateIncident(r)
	if err != nil {
		return aie, err
	}
	if err := PopulateIncidentStatus(aie); err != nil {
		return aie, err
	}
	return aie, nil
}

func PopulateIncident(url string) ([]db.Incident, error) {
	incidents := []db.Incident{}
	incident := db.Incident{}
	list, err := GetMainTable(url)
	if err != nil {
		return incidents, err
	}

	for index := 0; index < len(list); index += 6 {
		incident.IncidentNo = list[index]
		incident.IncidentType = list[index+1]
		if list[index+2] == "\u00a0\u00a0\u00a0\u00a0\u00a0\u00a0\u00a0" {
			incident.IncidentSubTupe = list[index+2]
			incident.Location = list[index+3]
			incident.Municipality = list[index+4]
		} else {
			incident.Location = list[index+2]
			incident.Municipality = list[index+3]
			incident.DispatchTime = list[index+4]
			incident.Station = list[index+5]
		}
		incidents = append(incidents, incident)
	}

	return incidents, nil
}

func PopulateIncidentStatus(aie *db.ActiveIncidentEntry) error {

	url := constants.WebCadURL + "livecad.asp?print=yes"
	r, err := Get(url)
	if err != nil {
		return err
	}

	_, incident, err := Tag(r)
	if err != nil {
		return err
	}
	for index, l := range incident {
		r, err = Get(GetDetail(l))
		if err != nil {
			return err
		}
		aie.IncidentWebPages = append(aie.IncidentWebPages, db.IncidentWebPage{Page: string(r)})
		if status, err := GetTable(r); err == nil {

			for i := 0; i < len(status); i += 3 {

				if len(status) <= i+3 {
					continue
				}

				aie.Incidents[index].IncidentStatus = append(
					aie.Incidents[index].IncidentStatus,
					db.IncidentStatus{TimeStamp: status[i],
						Unit:   status[i+1],
						Status: status[i+2]})

			}
		}
	}

	return nil
}
