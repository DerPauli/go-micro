package client

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/DerPauli/utils/utils"
	"github.com/buger/jsonparser"
	"github.com/joho/godotenv"
)

const base string = "https://api.twitter.com/2/"

var httpClient *http.Client
var token string

type Header struct {
	key   string
	value string
}

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	token = os.Getenv("TWITTER_TOKEN")
	if token == "" {
		log.Fatal("You must set your 'TWITTER_TOKEN' environmental variable.")
	}

	httpClient = &http.Client{}
}

func SetHeaders(rq *http.Request, headers []Header) *http.Request {
	for i := 0; i < len(headers); i++ {
		header := headers[i]
		rq.Header.Add(header.key, header.value)
	}

	return rq
}

func GetLatestTweets(from string) *utils.Tweet {
	query := "query=from:" + from + "&tweet.fields=created_at,public_metrics,source&expansions=author_id&user.fields=created_at"

	req, err := http.NewRequest("GET", base+"tweets/search/recent?"+query, nil)
	if err != nil {
		panic("Error while creating NewRequest")
	}

	tokenHeader := &Header{
		key:   "Authorization",
		value: "Bearer " + token,
	}
	headers := []Header{*tokenHeader}

	req = SetHeaders(req, headers)

	resp, err := httpClient.Do(req)
	if err != nil {
		panic("Error making request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic("HELP")
	}

	var tweets []*utils.Tweet

	jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		id, _, _, _ := jsonparser.Get(value, "id")
		authorId, _, _, _ := jsonparser.Get(value, "author_id")
		source, _, _, _ := jsonparser.Get(value, "source")
		createdAt, _, _, _ := jsonparser.Get(value, "created_at")
		text, _, _, _ := jsonparser.Get(value, "text")

		likeCount, _, _, _ := jsonparser.Get(value, "public_metrics", "like_count")
		retweetCount, _, _, _ := jsonparser.Get(value, "public_metrics", "retweet_count")
		replyCount, _, _, _ := jsonparser.Get(value, "public_metrics", "reply_count")
		quoteCount, _, _, _ := jsonparser.Get(value, "public_metrics", "quote_count")

		t := *utils.Tweet{
			id:           string(id),
			authorId:     string(authorId),
			source:       string(source),
			createdAt:    string(createdAt),
			text:         string(text),
			likeCount:    string(likeCount),
			retweetCount: string(retweetCount),
			replyCount:   string(replyCount),
			quoteCount:   string(quoteCount),
		}

		tweets = append(tweets, t)
	}, "data")

	return tweets
}
