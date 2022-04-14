package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/DerPauli/utils/utils"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}
}

func LoadTweets() {
	coll := client.Database("local").Collection("tweets")
	title := "Slash"

	var result bson.M
	cursor, err := coll.Find(context.TODO(), bson.D{})

	if err == mongo.ErrNoDocuments {
		log.Printf("No document was found with the title %s\n", title)
		return
	}

	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}

	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	log.Printf("%s\n", jsonData)

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}

func WriteTweets(tweets []*utils.Tweet) {

	for i := 0; i < len(tweets); i++ {
		log.Println(tweets[i].text)
	}
}
