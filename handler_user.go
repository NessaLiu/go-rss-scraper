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
func (apiConfig *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params) // decode into an instance of the params struct
	if err != nil {
		respondWithError(w, 400, fmt.Sprint("Error parsing JSON:", err)) // Sprint formats to string
		return
	}

	// Use DB to create user
	user, err := apiConfig.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	}) // sqlc created this method for us from reading our sql
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not create user: %s", err)) // Sprint formats to string
		return
	}
	respondWithJSON(w, 201, dbUserToUser(user))
}

func (apiConfig *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 200, dbUserToUser(user))
}

func (apiConfig *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiConfig.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get posts: %v", err))
		return
	}
	respondWithJSON(w, 200, dbPostsToPosts(posts))
}
