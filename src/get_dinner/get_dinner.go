package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// Global variables
var apiKey string
var apiUrl string
var sqsName string
var awsConfig aws.Config

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

type SubWeatherData struct {
	Weather []Weather `json:"weather"`
	Main    Main      `json:"main"`
	Wind    Wind      `json:"wind"`
	Clouds  Clouds    `json:"clouds"`
	Name    string    `json:"name"`
}

type GPTResponse struct {
	Meal                string
	CountryOfOrigin     string
	Ingredients         []Ingredient
	CookingInstructions []string
}

type Ingredient struct {
	Name   string `json:"name"`
	Amount string `json:"amount"`
}

// +------------------------+
// | Lambda entry functions |
// +------------------------+

// Lambda runner/worker.
func lambdaHandler(ctx context.Context, sqsEvent events.SQSEvent) (string, error) {

	weatherDataString := sqsEvent.Records[0].Body

	fmt.Println(weatherDataString)

	// Extract `sub-details` for ChatGPT to supply ideas
	/*
		subWeather := SubWeatherData{
			Weather: responseData.Weather,
			Main:    responseData.Main,
			Wind:    responseData.Wind,
			Clouds:  responseData.Clouds,
			Name:    responseData.Name,
		}
	*/

	// Place details on SQS queue for lambda processing.
	// fmt.Printf("%+v\n", subWeather)

	// sendWeatherToSqs(awsConfig, subWeather)

	// Inform operation is done.
	// return (fmt.Sprintf("City: %s", responseData.Name)), nil
	return "In development.", nil
}

// Lambda entry point
func main() {
	lambda.Start(lambdaHandler)
}

// +-----------+
// | Functions |
// +-----------+ 

// Retrieves and sets the Open weather map api key from AWS Parameter Store.
func getChatGPTApiKeyFromParameterStore(config aws.Config, secretKey string) {
	fmt.Println("Retrieving ChatGPT api-key from Parameter Store.")

	ssmClient := ssm.NewFromConfig(config)
	getApiKeyValue := &ssm.GetParameterInput{
		Name:           aws.String(secretKey),
		WithDecryption: aws.Bool(true),
	}

	result, err := ssmClient.GetParameter(context.TODO(), getApiKeyValue)
	if err != nil {
		fmt.Println("Retrieving ChatGPT api-key from Parameter Store failed.")
		fmt.Println(err)
		fmt.Println("Exiting...")
		os.Exit(1)
	}

	// Get the secret from the returned string.
	apiKey = *result.Parameter.Value
}

// Retrieves and sets the Open weather map URL from AWS Parameter Store.
func getChatGPTUrlFromParameterStore(config aws.Config, weatherUrlKey string) {
	fmt.Println("Retrieving ChatGPT url from Parameter Store.")

	ssmClient := ssm.NewFromConfig(config)
	getUrlValue := &ssm.GetParameterInput{
		Name:           aws.String(weatherUrlKey),
		WithDecryption: aws.Bool(false),
	}

	result, err := ssmClient.GetParameter(context.TODO(), getUrlValue)
	if err != nil {
		fmt.Println("Retrieving ChatGPT URL from Parameter store failed.")
		fmt.Println(err)
		fmt.Println("Exiting...")
		os.Exit(1)
	}

	// Get the secret from the returned string.
	apiUrl = *result.Parameter.Value
}

// Makes a Http request to the Open Weather Map API.
func makeHttpRequest(apiUrl string, apiKey string, city string) (string, error) {
	fmt.Println("Making Http request to Open weather map.")

	// Format the correct request URL and Send request.
	requestUrl := fmt.Sprintf(apiUrl, city, apiKey)
	response, err := http.Get(requestUrl)

	if err != nil {
		fmt.Println("Making request to open weather map failed.")
		fmt.Println(err)
		fmt.Println("Exiting...")
		os.Exit(1)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("WARNING. Non-200 HTTP response.")
		fmt.Printf("Response code: %v\n", response.StatusCode)
		fmt.Printf("Response reason: %s\n", response.Status)
		fmt.Println("Exiting...")
		os.Exit(1)
	}

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Reading response body failed.")
		fmt.Println(err)
		fmt.Println("Exiting...")
		os.Exit(1)
	}

	fmt.Println("Reading response body succeeded.")

	// Deserialize/decode from json into struct
	var weatherData string
	unmarshalError := json.Unmarshal([]byte(body), &weatherData)

	if unmarshalError != nil {
		fmt.Println("Deserializing Open Weather Map response failed.")
		fmt.Println(unmarshalError)
		fmt.Println("Exiting...")
		os.Exit(1)
	}

	fmt.Println("Deserializing Open Weather Map response succeeded.")

	return weatherData, nil
}

// Place Http result on Sqs.
func sendWeatherToSqs(config aws.Config, subWeather SubWeatherData) {
	fmt.Println("Publishing to Sqs.")

	sqsClient := sqs.NewFromConfig(config)
	queueInput := &sqs.GetQueueUrlInput{
		QueueName: aws.String(sqsName),
	}

	result, err := sqsClient.GetQueueUrl(context.TODO(), queueInput)

	if err != nil {
		fmt.Println("Retrieving the queue URL failed.")
		fmt.Println(err)
		fmt.Println("Exiting...")
		os.Exit(1)
	}

	queueUrl := result.QueueUrl

	// Serialize the sub-weather into JSON.
	jsonWeather, err := json.Marshal(subWeather)

	if err != nil {
		fmt.Println("Serializing weather data failed.")
		fmt.Println(err)
		fmt.Println("Exiting...")
		os.Exit(1)
	}

	sqsMessage := &sqs.SendMessageInput{
		MessageBody: aws.String(string(jsonWeather)),
		QueueUrl:    queueUrl,
	}

	_, errs := sqsClient.SendMessage(context.TODO(), sqsMessage)

	if errs != nil {
		fmt.Println("Publishing to Sqs failed.")
		fmt.Println(errs)
		fmt.Println("Exiting...")
		os.Exit(1)
	}

	fmt.Println("Publishing to Sqs sucessful.")
}
