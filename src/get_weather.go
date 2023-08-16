package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func LambdaHandler() (string, error) {
	return fmt.Sprintf("Hello, %v", "From AWS lambda."), nil
}

func main() {
	lambda.Start(LambdaHandler)
}
