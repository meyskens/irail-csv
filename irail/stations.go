package irail

type stationList struct {
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
	Station   []struct {
		LocationX    string `json:"locationX"`
		LocationY    string `json:"locationY"`
		ID           string `json:"id"`
		Standardname string `json:"standardname"`
		Name         string `json:"name"`
	} `json:"station"`
}

type Station struct {
	Name string `json:"name" csv:"name"`
}

// GetStationList gets all station names listed
func GetStationList() ([]Station, error) {
	resp := stationList{}

	r := getClient()
	r.SetResult(&resp)
	_, err := r.Get("/stations")

	if err != nil {
		return nil, err
	}

	names := []Station{}
	for _, station := range resp.Station {
		names = append(names, Station{Name: station.Name})
	}

	return names, nil
}
