package openweather

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenWeatherClient struct {
	apikey string
}

func New(apikey string) *OpenWeatherClient {
	return &OpenWeatherClient{
		apikey: apikey,
	}
}

func (o OpenWeatherClient) Coordinates(city string) (error, Coordinates) {
	url := "http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=5&appid=%s"
	resp, err := http.Get(fmt.Sprintf(url, city, o.apikey))
	if err != nil {
		return fmt.Errorf("error get coordinates: %w", err), Coordinates{}
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error fail get coordinates: %d", resp.StatusCode), Coordinates{}
	}

	var coordinatesResponse []CoordinatesResponse
	err = json.NewDecoder(resp.Body).Decode(&coordinatesResponse)
	if err != nil {
		return fmt.Errorf("error unmarshal: %w", err), Coordinates{}
	}

	if len(coordinatesResponse) == 0 {
		return fmt.Errorf("error empty coordinates: %w", err), Coordinates{}
	}

	return nil, Coordinates{
		Lat: coordinatesResponse[0].Lat,
		Lon: coordinatesResponse[0].Lon,
	}
}

func (o OpenWeatherClient) Weather(lat, lon float64) (error, Weather) {
	url := "https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric"
	resp, err := http.Get(fmt.Sprintf(url, lat, lon, o.apikey))
	if err != nil {
		return fmt.Errorf("error get weather: %w", err), Weather{}
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error fail get weather: %d", resp.StatusCode), Weather{}
	}

	var weatherResponse WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weatherResponse)
	if err != nil {
		return fmt.Errorf("error unmarshal weather response %w", err), Weather{}
	}

	return nil, Weather{
		Temp: weatherResponse.Main.Temp,
	}
}
