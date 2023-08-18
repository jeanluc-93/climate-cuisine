package main

import (
	"context"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

type WeatherResponse struct {
	Temperature float64 `json:"temperature"`
	// Add other fields you need
}

type MyEvent struct {
	Name string `json:"name"`
}

func loadAPIKey() string {
	// Load environmental variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("Weather API key not found")
	}
	return apiKey
}

func makeAPIRequest(apiKey string) (*resty.Response, error) {

	// api.openweathermap.org/data/2.5/weather?q=London,uk&APPID=
	// url := "https://weatherapi-com.p.rapidapi.com/current.json?q=53.1%2C-0.13"
	//url := "api.openweathermap.org/data/2.5/weather?q=Cape Town&units=metric&APPID=" + apiKey

	//req, _ := http.NewRequest("GET", url, nil)

	//req.Header.Add("X-RapidAPI-Key", "SIGN-UP-FOR-KEY")
	// req.Header.Add("X-RapidAPI-Host", "weatherapi-com.p.rapidapi.com")

	//res, _ := http.DefaultClient.Do(req)

	// defer res.Body.Close()
	//body, _ := io.ReadAll(res.Body)

	//fmt.Println(res)
	//fmt.Println(string(body))

	client := resty.New()
	response, err := client.R().Get("https://api.openweathermap.org/data/2.5/weather?q=Cape Town&units=metric&APPID=" + apiKey)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func lambdaHandler(ctx context.Context) (*resty.Response, error) {

	apiKey := loadAPIKey()
	response, err := makeAPIRequest(apiKey)

	if err != nil {
		log.Fatal("Error making API request.")
	}

	return response, nil
}

func main() {
	// lambda.Start(lambdaHandler)
	apiKey := loadAPIKey()
	response, err := makeAPIRequest(apiKey)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(response.Body())
}
