package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/labstack/echo"
	"github.com/meyskens/irail-csv/irail"
)

var stationList []irail.Station

func main() {
	var err error
	stationList, err = irail.GetStationList()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.GET("/stations", getStations)
	e.GET("/connections", getConnections)
	e.Logger.Fatal(e.Start(":80"))
}

func getStations(c echo.Context) error {
	list := []irail.Station{}
	if search := c.QueryParam("search"); search != "" {
		search = strings.ToLower(search)
		for _, station := range stationList {
			if strings.Contains(strings.ToLower(station.Name), search) {
				list = append(list, station)
			}
		}
	} else {
		list = stationList
	}

	out, err := gocsv.MarshalString(list)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Error while marshhaling csv")
	}
	c.Response().Header().Set("Content-Type", "text/csv")
	return c.String(http.StatusOK, out)
}

func getConnections(c echo.Context) error {
	connections, err := irail.GetConnection(
		c.QueryParam("from"),
		c.QueryParam("to"),
		c.QueryParam("timesel"),
		c.QueryParam("time"),
		c.QueryParam("date"))

	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Error while getting connections")
	}

	out, err := gocsv.MarshalString(connections)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Error while marshhaling csv")
	}
	c.Response().Header().Set("Content-Type", "text/csv")
	return c.String(http.StatusOK, out)
}
