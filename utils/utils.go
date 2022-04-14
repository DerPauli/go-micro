package utils

import (
	"bytes"
	"encoding/json"
)

// GLOBAL STRUCTS
type Tweet struct {
	id        string
	authorId  string
	source    string
	createdAt string
	text      string

	likeCount    string
	retweetCount string
	replyCount   string
	quoteCount   string
}

func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}
