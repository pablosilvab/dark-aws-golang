package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type BodyRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type BodyResponse struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// BodyRequest will be used to take the json response from client and build it
	bodyRequest := BodyRequest{}

	// Unmarshal the json, return 404 if error
	err := json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: string(buildResponse("Petición errónea", -1)), StatusCode: 404}, nil
	}

	// We will build the BodyResponse and send it back in json form
	bodyResponse := BodyResponse{
		FirstName: bodyRequest.FirstName,
		LastName:  bodyRequest.LastName,
	}

	// Marshal the response into json bytes, if error return 404
	response, err := json.Marshal(&bodyResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	//Returning response with AWS Lambda Proxy Response
	return events.APIGatewayProxyResponse{Body: string(response), StatusCode: 200}, nil
}

func buildResponse(message string, code int) []byte {
	response := ErrorResponse{message, code}
	js, _ := json.Marshal(&response)
	return js
}

func main() {
	lambda.Start(Handler)
}
