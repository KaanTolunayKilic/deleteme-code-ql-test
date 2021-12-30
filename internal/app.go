package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"git.haw-hamburg.de/aci822/gabi/tweet-extractor/internal/config"
	"git.haw-hamburg.de/aci822/gabi/tweet-extractor/internal/db/mongo"
	"git.haw-hamburg.de/aci822/gabi/tweet-extractor/internal/db/postgres"
	"git.haw-hamburg.de/aci822/gabi/tweet-extractor/internal/extraction"
	"git.haw-hamburg.de/aci822/gabi/tweet-extractor/internal/graph"
	"git.haw-hamburg.de/aci822/gabi/tweet-extractor/internal/graph/generated"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/rs/cors"

	"go.uber.org/zap"
)

const port string = "8080"

type App struct {
	logger              *zap.Logger
	config              config.Configuration
	tweetExtractor      extraction.TweetExtractor
	pgDatabaseClient    *postgres.PrismaClient
	mongoDatabaseClient *mongo.MongoDatabase
}

func NewApp(logger *zap.Logger, config config.Configuration) App {
	return App{
		logger: logger,
		config: config,
	}
}

func (app *App) Initialize() {

	pgClient := postgres.NewClient()
	if err := pgClient.Prisma.Connect(); err != nil {
		panic(err)
	}
	app.pgDatabaseClient = pgClient
	app.mongoDatabaseClient = mongo.NewDatastore(app.config, app.logger)

	app.tweetExtractor = extraction.NewTweetExtractor(app.config, app.logger, app.mongoDatabaseClient)
	// TODO: Implement synchronization at startup so that
	// activated queries are read from db and their extraction is started.
	app.syncExtractor()

}

func (app *App) syncExtractor() {
	ctx := context.Background()
	queries, err := app.pgDatabaseClient.SearchQuery.FindMany(
		postgres.SearchQuery.Active.Equals(true),
	).Exec(ctx)
	if err != nil {
		panic(err)
	}
	var queriesInner []postgres.InnerSearchQuery
	for index, query := range queries {
		queriesInner = append(queriesInner, query.InnerSearchQuery)
	}
	app.tweetExtractor.StartAllExtractions(queriesInner)
}

func (app *App) Shutdown() {
	app.tweetExtractor.StopAllExtractions()
	if err := app.pgDatabaseClient.Disconnect(); err != nil {
		// TODO: logging
		panic(err)
	}
}

func (app *App) Start() {
	app.logger.Info("Starting server")

	router := chi.NewRouter()
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		DatabaseClient: app.pgDatabaseClient,
		TweetExtractor: &app.tweetExtractor,
		Logger:         app.logger,
	}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	app.logger.Info(fmt.Sprintf("connect to http://localhost:%s/ for GraphQL playground", port))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
