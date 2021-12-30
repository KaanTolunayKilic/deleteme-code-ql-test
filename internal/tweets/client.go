package tweets

import (
	"context"

	"git.haw-hamburg.de/aci822/gabi/tweet-extractor/internal/config"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func NewClient(config config.Configuration) *twitter.Client {
	oauthConfig := oauth1.NewConfig(config.ConsumerKey, config.ConsumerSecret)
	oauthToken := oauth1.NewToken(config.AccessToken, config.AccessSecret)

	httpClient := oauthConfig.Client(context.Background(), oauthToken)

	return twitter.NewClient(httpClient)
}
