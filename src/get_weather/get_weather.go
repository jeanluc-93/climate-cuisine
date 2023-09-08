package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	// "github.com/go-resty/resty/v2"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

// Open Weather Map returned JSON as structs
type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
}

type Clouds struct {
	All int `json:"all"`
}

type Sys struct {
	Type    int    `json:"type"`
	ID      int    `json:"id"`
	Country string `json:"country"`
	Sunrise int    `json:"sunrise"`
	Sunset  int    `json:"sunset"`
}

type WeatherData struct {
	Coord      Coord     `json:"coord"`
	Weather    []Weather `json:"weather"`
	Base       string    `json:"base"`
	Main       Main      `json:"main"`
	Visibility int       `json:"visibility"`
	Wind       Wind      `json:"wind"`
	Clouds     Clouds    `json:"clouds"`
	Dt         int       `json:"dt"`
	Sys        Sys       `json:"sys"`
	Timezone   int       `json:"timezone"`
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Cod        int       `json:"cod"`
}

// Functions
func loadAPIKey() string {
	// Load environmental variables from .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("Weather API key not found")
	}
	return apiKey
}

func makeAPIRequest(apiKey string) (string, error) {

	requestUrl := "https://api.openweathermap.org/data/2.5/weather?q=Cape Town&units=metric&APPID=" + apiKey

	// Send GET request
	response, err := http.Get(requestUrl)

	if err != nil {
		fmt.Printf("Error sending GET request: %v\n", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		os.Exit(1)
	}

	// Deserialize/decode from json into struct
	var weatherData WeatherData
	unmarshalError := json.Unmarshal([]byte(body), &weatherData)

	if unmarshalError != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Got response")
	fmt.Printf(string(body))

	return string(body), nil
}

func lambdaHandler(ctx context.Context) (string, error) {

	apiKey := loadAPIKey()
	_, err := makeAPIRequest(apiKey)

	if err != nil {
		log.Fatal("Error making API request.")
	}

	//return response, nil
	return "Hello from lambda!", nil
}

func main() {
	lambda.Start(lambdaHandler)
}
