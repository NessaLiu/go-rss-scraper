package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/NessaLiu/go-rss-scraper/internal/database"
	"github.com/google/uuid"
)

// The signature of the handler can't change, so we make this a method and pass in the apiConfig struct
func (apiConfig *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params) // decode into an instance of the params struct
	if err != nil {
		respondWithError(w, 400, fmt.Sprint("Error parsing JSON:", err)) // Sprint formats to string
		return
	}

	// Use DB to create feed
	feed, err := apiConfig.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not create feed: %s", err)) // Sprint formats to string
		return
	}
	respondWithJSON(w, 201, dbFeedToFeed(feed))
}
