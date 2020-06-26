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
	"gopkg.in/mgo.v2/bson"
)

type Character struct {
	FirstName string `json:"firstName" bson:"firstname,omitempty"`
	LastName  string `json:"lastName" bson:"lastname,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("GetAll Function....")

	db, err := DBConnect()

	if err != nil {
		return events.APIGatewayProxyResponse{Body: string(buildResponse("Error al conectarse a base de datos", -1)), StatusCode: http.StatusInternalServerError}, nil
	}
	collection := db.Collection("characters")

	var characters []*Character

	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var character Character
		err := cur.Decode(&character)
		if err != nil {
			log.Fatal(err)
		}
		characters = append(characters, &character)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	response, err := json.Marshal(characters)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: string(buildResponse("Error interno", -1)), StatusCode: http.StatusInternalServerError}, nil
	}

	return events.APIGatewayProxyResponse{Body: string(response), StatusCode: http.StatusOK}, nil
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
