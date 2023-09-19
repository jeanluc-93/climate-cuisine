package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// Global variables
var apiKey string
var apiUrl string

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

// Received event struct
type Event struct {
	City string `json:"city"`
}

// +------------------------+
// | Lambda entry functions |
// +------------------------+

// Initialize function that sets up the running environment.
func init() {
	// Load environment keys from environment variables
	secretKey := os.Getenv("SECRET_KEY")
	weatherUrlKey := os.Getenv("OPEN_WEATHER_URL_KEY")
	region := os.Getenv("REGION") // af-south-1
	// sqsName := os.Getenv("SQS_QUEUE_NAME")

	// Load the AWS profile config
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	// Load API key and URL for usage.
	getApiKeyFromSecretsManager(cfg, secretKey)
	getUrlFromParameterStore(cfg, weatherUrlKey)
}

// Lambda runner/worker.
func lambdaHandler(ctx context.Context /*, event Event*/) (string, error) {

	// Make Http request to get daily weather
	responseData, responseError := makeHttpRequest(apiUrl, apiKey, "Cape Town")

	if responseError != nil {
		log.Fatal("Error making API request.")
	}

	// Extract `sub-details` for ChatGPT to supply ideas
	subWeather := WeatherData{
		Weather: responseData.Weather,
		Main:    responseData.Main,
		Wind:    responseData.Wind,
		Clouds:  responseData.Clouds,
		Name:    responseData.Name,
	}

	// Place details on SQS queue for lambda processing.
	fmt.Println(subWeather)

	// Inform operation is done.
	return (fmt.Sprintf("City: %s\n", responseData.Name)), nil
}

// Lambda entry point
func main() {
	lambda.Start(lambdaHandler)
}

// +-----------+
// | Functions |
// +-----------+

// Retrieves and sets the Open weather map API key from AWS Secrets Manager.
func getApiKeyFromSecretsManager(config aws.Config, secretKey string) {
	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	getSecretValue := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretKey),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), getSecretValue)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	// Get the secret from the returned string.
	apiKey = *result.SecretString
}

// Retrieves and sets the Open weather map URL from AWS Parameter Store.
func getUrlFromParameterStore(config aws.Config, weatherUrlKey string) {
	ssmClient := ssm.NewFromConfig(config)

	getUrlValue := &ssm.GetParameterInput{
		Name:           aws.String(weatherUrlKey),
		WithDecryption: aws.Bool(false),
	}

	result, err := ssmClient.GetParameter(context.TODO(), getUrlValue)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	// Get the secret from the returned string.
	apiUrl = *result.Parameter.Value
}

// Makes a Http request to the Open Weather Map API.
func makeHttpRequest(apiUrl string, apiKey string, city string) (WeatherData, error) {
	// Format the correct request URL.
	requestUrl := fmt.Sprintf(apiUrl, city, apiKey)

	// Send GET request
	response, err := http.Get(requestUrl)

	if err != nil {
		fmt.Printf("Error sending GET request: %v\n", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Not a 200 okay response.\n")
		fmt.Printf("Response code: %v\n", response.StatusCode)
		fmt.Printf("Response reason: %s\n", response.Status)

		os.Exit(1)
	}

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: \n %s \n", err)

		os.Exit(1)
	}

	// Deserialize/decode from json into struct
	var weatherData WeatherData
	unmarshalError := json.Unmarshal([]byte(body), &weatherData)

	if unmarshalError != nil {
		fmt.Println("Not able to deserialize Open Weather Map json response.")
		fmt.Printf("Error: \n %s \n", unmarshalError)

		os.Exit(1)
	}

	return weatherData, nil
}

func sendWeatherToSqs(config aws.Config, queue string, topic string, subWeather WeatherData) {
	sqsClient := sqs.NewFromConfig(config)

	// Get URL of queue
	queueInput := &sqs.GetQueueUrlInput{
		QueueName: aws.String("TestQueue"),
	}

	result, err := sqsClient.GetQueueUrl(context.TODO(), queueInput)

	if err != nil {
		fmt.Println("Got an error getting the queue URL:")
		fmt.Println(err)
		os.Exit(1)
	}

	queueUrl := result.QueueUrl

	// Serialize the sub-weather into JSON.
	jsonWeather, err := json.Marshal(subWeather)

	if err != nil {
		fmt.Println("Not able to serialize weather data to json.")
		os.Exit(1)
	}

	sqsMessage := &sqs.SendMessageInput{
		MessageBody: aws.String(string(jsonWeather)),
		QueueUrl:    queueUrl,
	}

	_, errs := sqsClient.SendMessage(context.TODO(), sqsMessage)

	if errs != nil {
		fmt.Println("Got an error sending the message:")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Weather message added to SQS.")
}
