package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type BodyRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bodyRequest := BodyRequest{}

	err := json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: string(buildResponse("Petición errónea", -1)), StatusCode: 404}, nil
	}

	if err != nil {
		return events.APIGatewayProxyResponse{Body: string(buildResponse("Error al ingresar registro", -1)), StatusCode: 500}, nil
	}

	response, err := json.Marshal(bodyRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: string(buildResponse("Error interno", -1)), StatusCode: 500}, nil
	}
	return events.APIGatewayProxyResponse{Body: string(response), StatusCode: 201}, nil
}

func buildResponse(message string, code int) []byte {
	log.Println("Error: ", message)
	response := ErrorResponse{message, code}
	js, _ := json.Marshal(&response)
	return js
}

func main() {
	lambda.Start(Handler)
}
