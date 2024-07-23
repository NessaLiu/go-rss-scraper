package main

import (
	"fmt"
	"net/http"

	"github.com/NessaLiu/go-rss-scraper/internal/auth"
	"github.com/NessaLiu/go-rss-scraper/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

// We want to create middleware where we have access to the user but can match the handler signature
func (apiConfig *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Authentication error: %v", err))
			return
		}
		user, err := apiConfig.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Failed to retrieve user: %v", err))
		}
		handler(w, r, user)
	}
}
