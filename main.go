package main

import (
	"github.com/DerPauli/go-micro/db/database"
	"github.com/DerPauli/go-micro/net/client"
)

func main() {
	database.Init()
	//database.LoadBooks()

	client.Init()
	tweets := client.GetLatestTweets("Tesla")
	database.WriteTweets(tweets)
}
