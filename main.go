package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
		FeelslikeC float64 `json:"feelslike_c"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int     `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				WillItRain   int `json:"will_it_rain"`
				ChanceOfRain int `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	q := "Ghaziabad"

	if len(os.Args) >= 2 {
		q = os.Args[1]
	}

	resp, err := http.Get("https://api.weatherapi.com/v1/forecast.json?key=de2687252e9e4f9f8dc131116242205&q=" + q + "&days=1&aqi=no&alerts=no")
	if err != nil {
		log.Fatalf("Error sending get request: %v", err)
	}
	defer resp.Body.Close()

	var weather Weather
	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		log.Fatalf("Error decoding response body: %v", err)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour
	fmt.Printf("%s, %s: %.0fC, %s\n", location.Name, location.Country, current.TempC, current.Condition.Text)

	for _, hour := range hours {
		date := time.Unix(int64(hour.TimeEpoch), 0)

		if date.Before(time.Now()) {
			continue
		}

		forecast := fmt.Sprintf("%s - %.0fC, %d, %s\n",
			date.Format("15:04"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text)

		if hour.ChanceOfRain == 100 {
			color.Blue(forecast)
		} else {
			fmt.Print(forecast)
		}
	}
}
