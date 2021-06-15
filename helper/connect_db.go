package helper

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const (
	CONNECTION_STRING = "mongodb://root:auth!@mongodb:27017"
	DB                = "db_authentication"
	USERS             = "users"
)

func ConnectDB() *mongo.Database {

	clientOptions := options.Client().ApplyURI(CONNECTION_STRING)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	database := client.Database(DB)
	//collection := client.Database(DB).Collection("users")

	return database
}

func HashPassword(password string, salt string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password+salt), 14)

	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}
