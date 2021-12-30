package tweets

import (
	"fmt"
	"log"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExtractedTweet struct {
	ID               primitive.ObjectID
	OriginalTweetID  int64
	CreatedAt        time.Time
	ExtractedAt      time.Time
	LastModified     time.Time
	Text             string
	IncludedHashtags []string
	Place            string
	FromUser         int64
	ExtractionQuery  []string
}

func getDefaultPlace() string {
	return ""
}

func ExtractTweet(tweet *twitter.Tweet, queries []string) *ExtractedTweet {
	createdAt, err := tweet.CreatedAtTime()
	if err != nil {
		fmt.Println(err.Error())
	}
	var hashtags []string
	for _, hashtag := range tweet.Entities.Hashtags {
		hashtags = append(hashtags, hashtag.Text)
	}
	place := getDefaultPlace()
	if len(place) > 2 {
		log.Printf("Place has forbidden length")
	}
	if tweet.Place != nil {
		place = tweet.Place.FullName
	} else {
		place = "undefined"
	}
	now := time.Now()
	return &ExtractedTweet{
		OriginalTweetID:  tweet.ID,
		CreatedAt:        createdAt,
		ExtractedAt:      now,
		LastModified:     now,
		Text:             tweet.Text,
		IncludedHashtags: hashtags,
		Place:            place,
		FromUser:         tweet.User.ID,
		ExtractionQuery:  queries,
	}
}
