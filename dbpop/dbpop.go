package dbpop

import (
	

	"github.com/cwxstat/activeIncident/constants"
	"github.com/cwxstat/activeIncident/db"
	"github.com/cwxstat/activeIncident/dbutils"
	"github.com/cwxstat/activeIncident/scrape"
)

// PopulateActiveIncidentEntry populates the ActiveIncidentEntry from a web
func PopulateActiveIncidentEntry() (*db.ActiveIncidentEntry, error) {

	aie := &db.ActiveIncidentEntry{}
	url := constants.WebCadMontcoPrint
	r, err := scrape.Get(url)
	if err != nil {
		return aie, err
	}
	aie.MainWebPage = r
	aie.TimeStamp = dbutils.NYtime()
	aie.Incidents, err = PopulateIncident(r)
	if err != nil {
		return aie, err
	}
	if err := PopulateIncidentStatus(aie); err != nil {
		return aie, err
	}
	return aie, nil
}

// PopulateIncident populates the Incident from a web.
func PopulateIncident(url string) ([]db.Incident, error) {
	incidents := []db.Incident{}
	incident := db.Incident{}
	list, err := scrape.GetMainTable(url)
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

// PopulateIncidentStatus populates the IncidentStatus from a web. Status
// is the status of the incident, which is "Enroute", "Dispatched", "Arrived" ...
func PopulateIncidentStatus(aie *db.ActiveIncidentEntry) error {

	url := constants.WebCadMontcoPrint
	r, err := scrape.Get(url)
	if err != nil {
		return err
	}

	_, incident, err := scrape.Tag(r)
	if err != nil {
		return err
	}
	for index, l := range incident {
		r, err = scrape.Get(scrape.GetDetail(l))
		if err != nil {
			return err
		}
		aie.IncidentWebPages = append(aie.IncidentWebPages, db.IncidentWebPage{Page: string(r)})
		if status, err := scrape.GetTable(r); err == nil {

			for i := 0; i < len(status); i += 3 {

				if len(status) < i+3 {
					continue
				}

				if len(status) == i+2 {
					aie.Incidents[index].IncidentStatus = append(
						aie.Incidents[index].IncidentStatus,
						db.IncidentStatus{TimeStamp: status[i],
							Status: status[i+1]})
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
