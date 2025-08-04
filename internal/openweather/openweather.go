// Package openweather is a nice package
package openweather

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type OpenWeatherClient struct {
	APIKey string
}

func New(APIKey string) *OpenWeatherClient {
	return &OpenWeatherClient{
		APIKey: APIKey,
	}
}

func (o OpenWeatherClient) GetCoordinates(city string) (Coordinates, error) {
	const op = "openweather.GetCoordinates"

	url := "http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=5&appid=%s"
	resp, err := http.Get(fmt.Sprintf(url, city, o.APIKey))
	if err != nil {
		return Coordinates{}, fmt.Errorf("%s: failed to get coornidates: %w", op, err)
	}

	if resp.StatusCode != 200 {
		return Coordinates{}, fmt.Errorf("%s: status code: %d", op, resp.StatusCode)
	}

	var coordinatesResponse []CoordinateResponse
	if err := json.NewDecoder(resp.Body).Decode(&coordinatesResponse); err != nil {
		return Coordinates{}, fmt.Errorf("%s: failed to decode response body: %w", op, err)
	}

	if len(coordinatesResponse) == 0 {
		return Coordinates{}, fmt.Errorf("%s: empty response", op)
	}

	name := coordinatesResponse[0].Name

	if coordinatesResponse[0].LocalNames["ru"] != "" {
		name = coordinatesResponse[0].LocalNames["ru"]
	}

	return Coordinates{
		Name: name,
		Lat: coordinatesResponse[0].Latitude,
		Lon: coordinatesResponse[0].Longitude,
	}, nil
}
func (o OpenWeatherClient) GetWeather(name string, lat, lon float64) (Weather, error) {
	const op = "openweather.GetWeather"

	url := "https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric"

	resp, err := http.Get(fmt.Sprintf(url, lat, lon, o.APIKey))
	if err != nil {
		return Weather{}, fmt.Errorf("%s: failed with GET request: %w", op, err)
	}
	if resp.StatusCode != 200 {
		return Weather{}, fmt.Errorf("%s: status code %d", op, resp.StatusCode)
	}

	var weatherResponse WeatherResponse

	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return Weather{}, fmt.Errorf("%s: failed to decode response body: %w", op, err)
	}

	log.Println(weatherResponse.Weather)
	
	return Weather{
		Weather: weatherResponse.Weather[0].Main,
		Description: weatherResponse.Weather[0].Description,
		Temp: weatherResponse.Main.Temp,
		FeelsLike: weatherResponse.Main.FeelsLike,
		Speed: weatherResponse.Wind.Speed,
	}, nil
}
