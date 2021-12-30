package extraction

import (
	"context"

	"git.haw-hamburg.de/aci822/gabi/tweet-extractor/internal/config"
	"git.haw-hamburg.de/aci822/gabi/tweet-extractor/internal/db/mongo"
	"git.haw-hamburg.de/aci822/gabi/tweet-extractor/internal/db/postgres"
	"git.haw-hamburg.de/aci822/gabi/tweet-extractor/internal/tweets"
	"github.com/dghubble/go-twitter/twitter"
	"go.uber.org/zap"
)

type TweetExtractor struct {
	QueryExtractionStreams map[string]chan string
	Configuration          config.Configuration
	Logger                 *zap.Logger
	mongoDatabaseClient    *mongo.MongoDatabase
}

func NewTweetExtractor(config config.Configuration, logger *zap.Logger, mongoDatabaseClient *mongo.MongoDatabase) TweetExtractor {
	queryExtractionStreams := make(map[string](chan string))
	return TweetExtractor{
		QueryExtractionStreams: queryExtractionStreams,
		Configuration:          config,
		Logger:                 logger,
		mongoDatabaseClient:    mongoDatabaseClient,
	}
}

func (te *TweetExtractor) StartAllExtractions(queries []postgres.InnerSearchQuery) {
	for _, query := range queries {
		te.StartExtraction(query)
	}
}

func (te *TweetExtractor) StartExtraction(query postgres.InnerSearchQuery) {
	client := tweets.NewClient(te.Configuration)
	te.QueryExtractionStreams[query.ID] = te.extract(client, query, te.Logger)
}

func (te *TweetExtractor) StopAllExtractions() {
	for _, val := range te.QueryExtractionStreams {
		val <- "exit"
	}
}

func (te *TweetExtractor) StopExtraction(queryId string) {
	te.QueryExtractionStreams[queryId] <- "exit"
	delete(te.QueryExtractionStreams, queryId)
}

func (te *TweetExtractor) extract(client *twitter.Client, query postgres.InnerSearchQuery, logger *zap.Logger) chan string {
	logger.Info("Starting extraction.", zap.String("query", query.ID))
	ctx := context.Background()
	exit := make(chan string)
	go func() {
		params := &twitter.StreamFilterParams{
			Language:      []string{"de"},
			Track:         query.Tags,
			StallWarnings: twitter.Bool(true),
		}

		stream, err := client.Streams.Filter(params)

		if err != nil {
			logger.Error("Error while creating Twitter Stream.", zap.String("originalError", err.Error()))
		}

		demux := twitter.NewSwitchDemux()
		demux.Tweet = func(tweet *twitter.Tweet) {
			te.Logger.Info("Found tweet")
			extractedTweet := tweets.ExtractTweet(tweet, query.Tags)
			_, err := te.mongoDatabaseClient.Session.Database(te.Configuration.Mongo.Database).Collection(te.Configuration.Mongo.Collection).InsertOne(ctx, extractedTweet)
			if err != nil {
				logger.Error("Error saving extracted tweet to database", zap.String("originalError", err.Error()))
			}
		}

		for message := range stream.Messages {
			select {
			case <-exit:
				logger.Info("Stoping extraction.", zap.String("query", query.ID))
				stream.Stop()
				return
			default:
				demux.Handle(message)
			}
		}
	}()
	return exit
}
