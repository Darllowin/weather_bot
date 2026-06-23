package openweather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type OpenWeatherClient struct {
	apikey string
	client *http.Client
}

func New(apikey string) *OpenWeatherClient {
	return &OpenWeatherClient{
		apikey: apikey,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (o OpenWeatherClient) Coordinates(city string) (Coordinates, error) {
	url := "http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=5&appid=%s"
	resp, err := o.client.Get(fmt.Sprintf(url, city, o.apikey))
	if err != nil {
		return Coordinates{}, fmt.Errorf("error get coordinates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Coordinates{}, fmt.Errorf("error fail get coordinates: %d", resp.StatusCode)
	}

	var coordinatesResponse []CoordinatesResponse
	err = json.NewDecoder(resp.Body).Decode(&coordinatesResponse)
	if err != nil {
		return Coordinates{}, fmt.Errorf("error unmarshal: %w", err)
	}

	if len(coordinatesResponse) == 0 {
		return Coordinates{}, fmt.Errorf("empty coordinates response for city: %s", city)
	}

	return Coordinates{
		Lat: coordinatesResponse[0].Lat,
		Lon: coordinatesResponse[0].Lon,
	}, nil
}

func (o OpenWeatherClient) Weather(lat, lon float64) (Weather, error) {
	url := "https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric"
	resp, err := o.client.Get(fmt.Sprintf(url, lat, lon, o.apikey))
	if err != nil {
		return Weather{}, fmt.Errorf("error get weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Weather{}, fmt.Errorf("error fail get weather: %d", resp.StatusCode)
	}

	var weatherResponse WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weatherResponse)
	if err != nil {
		return Weather{}, fmt.Errorf("error unmarshal weather response %w", err)
	}

	return Weather{
		Temp: weatherResponse.Main.Temp,
	}, nil
}
