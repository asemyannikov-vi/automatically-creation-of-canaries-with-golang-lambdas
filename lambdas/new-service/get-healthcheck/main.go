package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func healthCheck() (*events.APIGatewayProxyResponse, error) {
	response, err := http.Get("https://dev.cloud.virtana.com/uin/api/v1/health")
	if err != nil {
		fmt.Println("Failed to send GET request. Error:", err.Error())
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read body of response. Error:", err.Error())
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: response.StatusCode,
	}, nil
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	response, err := healthCheck()
	if err != nil {
		fmt.Println("Failed to handle healthCheck endpoint. Error:", err.Error())
		return nil, err
	}
	return response, nil
}

func main() {
	lambda.Start(handler)
}
