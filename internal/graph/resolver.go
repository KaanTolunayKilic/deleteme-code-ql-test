package graph

import (
	"git.haw-hamburg.de/aci822/gabi/tweet-extractor/internal/db/postgres"
	"git.haw-hamburg.de/aci822/gabi/tweet-extractor/internal/extraction"
	"go.uber.org/zap"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DatabaseClient *postgres.PrismaClient
	TweetExtractor *extraction.TweetExtractor
	Logger         *zap.Logger
}
