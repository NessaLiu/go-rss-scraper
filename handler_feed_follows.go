package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/NessaLiu/go-rss-scraper/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiConfig *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params) // decode into an instance of the params struct
	if err != nil {
		respondWithError(w, 400, fmt.Sprint("Error parsing JSON:", err)) // Sprint formats to string
		return
	}

	feedFollow, err := apiConfig.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not create feed follow: %s", err)) // Sprint formats to string
		return
	}
	respondWithJSON(w, 201, dbFeedFollowToFeedFollow(feedFollow))
}

func (apiConfig *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := apiConfig.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get feed follows for user: %s", err)) // Sprint formats to string
		return
	}
	respondWithJSON(w, 201, dbFeedFollowsToFeedFollows(feedFollows))
}

func (apiConfig *apiConfig) handlerDeleteFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIdStr := chi.URLParam(r, "feedFollowID") // feedFollowID is what matches the arg in the delete path in the {}
	feedFollowId, err := uuid.Parse(feedFollowIdStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not parse feed follow id: %s", err))
		return
	}
	err = apiConfig.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowId,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not delete feed follow: %s", err)) // Sprint formats to string
		return
	}
	respondWithJSON(w, 200, struct{}{}) // returning an empty json object, they probably just care about the 200
}
