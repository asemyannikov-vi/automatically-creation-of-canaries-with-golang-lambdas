package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func healthCheck() (*events.APIGatewayProxyResponse, error) {
	cookieUserId := &http.Cookie{
		Name:   "x-user-id",
		Value:  "eaa46a16-cf82-4621-a8c6-6eedce335e94",
		MaxAge: 300,
	}
	cookieIdToken := &http.Cookie{
		Name:   "x-id-token",
		Value:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0b206YWN0aXZlT3JnYW5pemF0aW9uSWQiOiJjMWIxYmVlMi00MTcxLTRmOWUtYjc3OC0wMDkzY2E2ZGY2ZDIiLCJjdXN0b206YWN0aXZlT3JnYW5pemF0aW9uTmFtZSI6IlZpcnRhbmEiLCJjdXN0b206YXNzdW1lcklkIjoiIiwiY3VzdG9tOmVtYWlsIjoiamRvd0B2aXJ0YW5hLmNvbSIsImN1c3RvbTpnbG9iYWxSb2xlcyI6W10sImN1c3RvbTpuYW1lIjoiSiBEIiwiY3VzdG9tOm9yZ2FuaXphdGlvbnMiOlt7ImlkIjoiYzFiMWJlZTItNDE3MS00ZjllLWI3NzgtMDA5M2NhNmRmNmQyIiwibmFtZSI6IlZpcnRhbmEiLCJyb2xlcyI6W3siaWQiOjEzNDcsIm5hbWUiOiJ2aXJ0YW5hLXBsYXRmb3JtLWFkbWluIn1dfV0sImN1c3RvbTpwYXJlbnRPcmdJZCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsImV4cCI6MTYzNzkxNzU4MiwiaWF0IjoxNjM3ODMxMTgyLCJpc3MiOiJhdXRoLXNlcnZpY2UiLCJuYmYiOjE2Mzc4MzExODIsInN1YiI6ImVhYTQ2YTE2LWNmODItNDYyMS1hOGM2LTZlZWRjZTMzNWU5NCJ9.AlkTX7xgXReXEgVeGfo2WaypkRcaEgF7X9fiUR4Mw8Y",
		MaxAge: 300,
	}

	request, err := http.NewRequest("POST", "https://dev.cloud.virtana.com/uin/api/v1/messages", nil)
	if err != nil {
		fmt.Println("Failed to create a POST request. Error:", err.Error())
		return nil, err
	}

	request.AddCookie(cookieUserId)
	request.AddCookie(cookieIdToken)

	var client http.Client
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to send a POST request. Error:", err.Error())
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
