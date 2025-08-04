// Package openweather is a nice package
package openweather

type CoordinateResponse struct {
	Name string `json:"name"`
	LocalNames map[string]string `json:"local_names"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

type Coordinates struct {
	Name string
	Lat  float64
	Lon  float64
}

type WeatherResponse struct {
	Weather []struct {
		Main string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
}

type Weather struct {
	Weather string
	Description string
	Temp      float64
	FeelsLike float64
	Speed float64
}
