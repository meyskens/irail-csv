package irail

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	os.Setenv("TZ", "Europe/Brussels")
}

type conectionList struct {
	Version    string `json:"version"`
	Timestamp  string `json:"timestamp"`
	Connection []struct {
		ID        string `json:"id"`
		Departure struct {
			Delay       string `json:"delay"`
			Station     string `json:"station"`
			Stationinfo struct {
				LocationX    string `json:"locationX"`
				LocationY    string `json:"locationY"`
				ID           string `json:"id"`
				Standardname string `json:"standardname"`
				Name         string `json:"name"`
			} `json:"stationinfo"`
			Time         string `json:"time"`
			Vehicle      string `json:"vehicle"`
			Platform     string `json:"platform"`
			Platforminfo struct {
				Name   string `json:"name"`
				Normal string `json:"normal"`
			} `json:"platforminfo"`
			Canceled            string `json:"canceled"`
			DepartureConnection string `json:"departureConnection"`
			Direction           struct {
				Name string `json:"name"`
			} `json:"direction"`
			Left      string `json:"left"`
			Walking   string `json:"walking"`
			Occupancy struct {
				ID   string `json:"@id"`
				Name string `json:"name"`
			} `json:"occupancy"`
		} `json:"departure"`
		Arrival struct {
			Delay       string `json:"delay"`
			Station     string `json:"station"`
			Stationinfo struct {
				LocationX    string `json:"locationX"`
				LocationY    string `json:"locationY"`
				ID           string `json:"id"`
				Standardname string `json:"standardname"`
				Name         string `json:"name"`
			} `json:"stationinfo"`
			Time         string `json:"time"`
			Vehicle      string `json:"vehicle"`
			Platform     string `json:"platform"`
			Platforminfo struct {
				Name   string `json:"name"`
				Normal string `json:"normal"`
			} `json:"platforminfo"`
			Canceled  string `json:"canceled"`
			Direction struct {
				Name string `json:"name"`
			} `json:"direction"`
			Arrived string `json:"arrived"`
			Walking string `json:"walking"`
		} `json:"arrival"`
		Duration string `json:"duration"`
		Vias     struct {
			Number string `json:"number"`
			Via    []struct {
				ID      string `json:"id"`
				Arrival struct {
					Time         string `json:"time"`
					Platform     string `json:"platform"`
					Platforminfo struct {
						Name   string `json:"name"`
						Normal string `json:"normal"`
					} `json:"platforminfo"`
					IsExtraStop string `json:"isExtraStop"`
					Delay       string `json:"delay"`
					Canceled    string `json:"canceled"`
					Arrived     string `json:"arrived"`
					Walking     string `json:"walking"`
					Direction   struct {
						Name string `json:"name"`
					} `json:"direction"`
					Vehicle             string `json:"vehicle"`
					DepartureConnection string `json:"departureConnection"`
				} `json:"arrival"`
				Departure struct {
					Time         string `json:"time"`
					Platform     string `json:"platform"`
					Platforminfo struct {
						Name   string `json:"name"`
						Normal string `json:"normal"`
					} `json:"platforminfo"`
					IsExtraStop string `json:"isExtraStop"`
					Delay       string `json:"delay"`
					Canceled    string `json:"canceled"`
					Left        string `json:"left"`
					Walking     string `json:"walking"`
					Direction   struct {
						Name string `json:"name"`
					} `json:"direction"`
					Vehicle             string `json:"vehicle"`
					DepartureConnection string `json:"departureConnection"`
					Occupancy           struct {
						ID   string `json:"@id"`
						Name string `json:"name"`
					} `json:"occupancy"`
				} `json:"departure"`
				TimeBetween string `json:"timeBetween"`
				Station     string `json:"station"`
				Stationinfo struct {
					LocationX    string `json:"locationX"`
					LocationY    string `json:"locationY"`
					ID           string `json:"id"`
					Standardname string `json:"standardname"`
					Name         string `json:"name"`
				} `json:"stationinfo"`
				Vehicle   string `json:"vehicle"`
				Direction struct {
					Name string `json:"name"`
				} `json:"direction"`
			} `json:"via"`
		} `json:"vias,omitempty"`
		Occupancy struct {
			ID   string `json:"@id"`
			Name string `json:"name"`
		} `json:"occupancy"`
	} `json:"connection"`
}

type Connection struct {
	DepartureTime string `json:"departureTime" csv:"departure_time"`
	ArrivalTime   string `json:"arrivalTime" csv:"arrival_time"`
	TrainTypes    string `json:"trainTypes" csv:"train_types"` // string of all taken train types eg "IC IC L"
	ViaText       string `json:"viaText" csv:"via_text"`       // string with instrutction of vias
}

// GetConnection gets a connection between 2 stations
func GetConnection(from, to, timesel, timeString, dateString string) ([]Connection, error) {
	resp := conectionList{}

	r := getClient()
	r.SetResult(&resp)
	r.SetQueryParam("typeOfTransport", "nointernationaltrains") // Let's keep it with NMBS
	r.SetQueryParam("results", "5")                             // Hardcoded for now

	r.SetQueryParam("from", from)
	r.SetQueryParam("to", to)
	if timesel != "" {
		r.SetQueryParam("timesel", timesel)
	}
	if timeString != "" {
		r.SetQueryParam("time", timeString)
	}
	if dateString != "" {
		r.SetQueryParam("date", dateString)
	}

	_, err := r.Get("/connections/")

	if err != nil {
		return nil, err
	}

	connections := []Connection{}
	for _, connection := range resp.Connection {
		depTime, _ := strconv.ParseInt(connection.Departure.Time, 10, 64)
		depHour, depMin, _ := time.Unix(depTime, 0).Clock()

		arrTime, _ := strconv.ParseInt(connection.Arrival.Time, 10, 64)
		arrHour, arrMin, _ := time.Unix(arrTime, 0).Clock()

		trainTypes := getVehicleType(connection.Departure.Vehicle) + " "
		viaText := ""
		if len(connection.Vias.Via) == 0 {
			viaText = fmt.Sprintf("%d:%d spoor %s: %s %s %d:%d",
				depHour, depMin,
				connection.Departure.Platforminfo.Name,
				getVehicleNumber(connection.Departure.Vehicle),
				connection.Departure.Direction.Name,
				arrHour, arrMin)
		}
		for i, via := range connection.Vias.Via {
			var depTime int64
			if i == 0 {
				depTime, _ = strconv.ParseInt(connection.Departure.Time, 10, 64)
			} else {
				depTime, _ = strconv.ParseInt(connection.Vias.Via[i-1].Departure.Time, 10, 64)
			}

			viaDepHour, viaDepMin, _ := time.Unix(depTime, 0).Clock()
			arrTime, _ := strconv.ParseInt(via.Arrival.Time, 10, 64)
			viaArrHour, viaArrMin, _ := time.Unix(arrTime, 0).Clock()

			trainTypes += getVehicleType(via.Vehicle) + " "
			viaText += fmt.Sprintf("%d:%d spoor %s: %s %s %d:%d | ",
				viaDepHour, viaDepMin,
				via.Departure.Platforminfo.Name,
				getVehicleNumber(via.Vehicle),
				via.Direction.Name,
				viaArrHour, viaArrMin)
		}

		if len(connection.Vias.Via) != 0 {
			depTime, _ = strconv.ParseInt(connection.Vias.Via[len(connection.Vias.Via)-1].Departure.Time, 10, 64)
			viaDepHour, viaDepMin, _ := time.Unix(depTime, 0).Clock()

			viaText += fmt.Sprintf("%d:%d spoor %s: %s %s %d:%d | ",
				viaDepHour, viaDepMin,
				connection.Vias.Via[len(connection.Vias.Via)-1].Departure.Platforminfo.Name,
				getVehicleNumber(connection.Arrival.Vehicle),
				connection.Arrival.Direction.Name,
				arrHour, arrMin)
		}

		trainTypes = strings.TrimRight(trainTypes, " ")
		viaText = strings.TrimRight(viaText, " | ")

		connections = append(connections, Connection{
			DepartureTime: fmt.Sprintf("%d:%d", depHour, depMin),
			ArrivalTime:   fmt.Sprintf("%d:%d", arrHour, arrMin),
			TrainTypes:    trainTypes,
			ViaText:       viaText,
		})
	}

	return connections, nil
}

func getVehicleType(in string) string {
	name := getVehicleNumber(in)
	re := regexp.MustCompile("[^A-Z]")
	return re.ReplaceAllString(name, "")
}

func getVehicleNumber(in string) string {
	parts := strings.Split(in, ".")
	return parts[len(parts)-1]
}
