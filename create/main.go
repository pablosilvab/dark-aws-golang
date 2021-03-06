package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CharacterRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	newCharacter := CharacterRequest{}
	log.Println("Create Function....")
	log.Println(request.Body)

	err := json.Unmarshal([]byte(request.Body), &newCharacter)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: string(buildResponse("Petición errónea", -1)), StatusCode: http.StatusBadRequest}, nil
	}

	db, err := DBConnect()

	if err != nil {
		return events.APIGatewayProxyResponse{Body: string(buildResponse("Error al conectarse a base de datos", -1)), StatusCode: http.StatusInternalServerError}, nil
	}
	collection := db.Collection("characters")
	result, err := collection.InsertOne(context.TODO(), newCharacter)

	if err != nil {
		return events.APIGatewayProxyResponse{Body: string(buildResponse("Error al ingresar registro", -1)), StatusCode: http.StatusInternalServerError}, nil
	}

	response, err := json.Marshal(result)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: string(buildResponse("Error interno", -1)), StatusCode: http.StatusInternalServerError}, nil
	}
	return events.APIGatewayProxyResponse{Body: string(response), StatusCode: http.StatusCreated}, nil
}

func buildResponse(message string, code int) []byte {
	log.Println("Error: ", message)
	response := ErrorResponse{message, code}
	js, _ := json.Marshal(&response)
	return js
}

func DBConnect() (*mongo.Database, error) {

	// Set client options
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Println("Error en DBConnect")
		return nil, err
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Println("Error Ping")
		return nil, err
	}

	log.Println("Connected to MongoDB!")
	return client.Database("dark"), nil

}

func main() {
	lambda.Start(Handler)
}
