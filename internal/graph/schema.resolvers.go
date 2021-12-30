package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	"git.haw-hamburg.de/aci822/gabi/tweet-extractor/internal/db/postgres"
	"git.haw-hamburg.de/aci822/gabi/tweet-extractor/internal/graph/generated"
	"git.haw-hamburg.de/aci822/gabi/tweet-extractor/internal/graph/model"
	"go.uber.org/zap"
)

func (r *mutationResolver) CreateTalkshow(ctx context.Context, newTalkshow model.NewTalkshow) (*postgres.InnerTalkshow, error) {
	createdTalkshow, err := r.DatabaseClient.Talkshow.CreateOne(
		postgres.Talkshow.Kanal.Set(newTalkshow.Kanal),
		postgres.Talkshow.Host.Set(newTalkshow.Host),
	).Exec(ctx)
	r.Logger.Info("created talkshow", zap.String("host", createdTalkshow.Host))
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("internal server error")
	}
	return &createdTalkshow.InnerTalkshow, nil
}

func (r *mutationResolver) AddQueryToTalkshow(ctx context.Context, newSearchQuery model.NewSearchQuery, talkshow string) (*postgres.InnerTalkshow, error) {
	query, err := r.DatabaseClient.SearchQuery.CreateOne(
		postgres.SearchQuery.Active.Set(newSearchQuery.Active),
		postgres.SearchQuery.Talkshow.Link(
			postgres.Talkshow.ID.Equals(talkshow),
		),
		postgres.SearchQuery.Tags.Set(newSearchQuery.Tags),
	).Exec(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("internal server error")
	}
	if query.Active {
		r.TweetExtractor.StartExtraction(query.InnerSearchQuery)
	}
	updatedTalkshow, err := r.DatabaseClient.Talkshow.FindUnique(
		postgres.Talkshow.ID.Equals(talkshow),
	).Exec(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("internal server error")
	}
	return &updatedTalkshow.InnerTalkshow, nil
}

func (r *mutationResolver) ToggleSearchQueryState(ctx context.Context, searchQueryID string, targetState bool) (*postgres.InnerSearchQuery, error) {
	updatedSearchQuery, err := r.DatabaseClient.SearchQuery.FindUnique(
		postgres.SearchQuery.ID.Equals(searchQueryID),
	).Update(
		postgres.SearchQuery.Active.Set(targetState),
	).Exec(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("internal server error")
	}
	if updatedSearchQuery.Active {
		r.TweetExtractor.StartExtraction(updatedSearchQuery.InnerSearchQuery)
	} else {
		r.TweetExtractor.StopExtraction(updatedSearchQuery.ID)
	}
	return &updatedSearchQuery.InnerSearchQuery, nil
}

func (r *queryResolver) Talkshows(ctx context.Context) ([]postgres.InnerTalkshow, error) {
	rawTalkshows, err := r.DatabaseClient.Talkshow.FindMany().Exec(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("internal server error")
	}
	var talkshows []postgres.InnerTalkshow
	for _, talkshow := range rawTalkshows {
		talkshows = append(talkshows, talkshow.InnerTalkshow)
	}
	return talkshows, nil
}

func (r *searchQueryResolver) Talkshow(ctx context.Context, obj *postgres.InnerSearchQuery) (*postgres.InnerTalkshow, error) {
	searchQueryWithTalkshow, err := r.DatabaseClient.SearchQuery.FindUnique(
		postgres.SearchQuery.ID.Equals(obj.ID),
	).With(
		postgres.SearchQuery.Talkshow.Fetch(),
	).Exec(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("internal server error")
	}
	return &searchQueryWithTalkshow.Talkshow().InnerTalkshow, nil
}

func (r *talkshowResolver) Queries(ctx context.Context, obj *postgres.InnerTalkshow) ([]postgres.InnerSearchQuery, error) {
	talkshowWithQueries, err := r.DatabaseClient.Talkshow.FindUnique(
		postgres.Talkshow.ID.Equals(obj.ID),
	).With(
		postgres.Talkshow.Queries.Fetch(),
	).Exec(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("internal server error")
	}
	var queries []postgres.InnerSearchQuery
	for _, query := range talkshowWithQueries.Queries() {
		queries = append(queries, query.InnerSearchQuery)
	}
	return queries, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// SearchQuery returns generated.SearchQueryResolver implementation.
func (r *Resolver) SearchQuery() generated.SearchQueryResolver { return &searchQueryResolver{r} }

// Talkshow returns generated.TalkshowResolver implementation.
func (r *Resolver) Talkshow() generated.TalkshowResolver { return &talkshowResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type searchQueryResolver struct{ *Resolver }
type talkshowResolver struct{ *Resolver }
