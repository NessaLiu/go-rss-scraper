package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/NessaLiu/go-rss-scraper/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	// feed, err := urlToFeed("https://wagslane.dev/index.xml")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(feed)

	godotenv.Load(".env")

	portStr := os.Getenv("PORT")
	if portStr == "" {
		log.Fatal("PORT is not found in the env")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the env")
	}

	// Connect to DB
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Unable to connect to database")
	}
	db := database.New(conn)
	apiConfig := apiConfig{
		DB: db,
	} // create api config

	go startScraping(db, 10, time.Minute) // start a new go routine to not disrupt the main flow (startScraping never returns!)

	router := chi.NewRouter()
	// cors configuration lets people send requests to our server from a browser
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()                // create v1 router so we can mount it to the v1 path
	v1Router.Get("/healthz", handlerReadiness) // Connecting handlerReadiness function to /healthz path (scope to GET requests)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiConfig.handlerCreateUser)
	// Calling middleware auth func to get authenticated user, then calling get user handler
	v1Router.Get("/users", apiConfig.middlewareAuth(apiConfig.handlerGetUser))

	v1Router.Post("/feeds", apiConfig.middlewareAuth(apiConfig.handlerCreateFeed))
	v1Router.Get("/feeds", apiConfig.handlerGetFeeds)

	v1Router.Post("/feed_follows", apiConfig.middlewareAuth(apiConfig.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiConfig.middlewareAuth(apiConfig.handlerGetFeedFollows))
	// Delete typically doesn't have a body for the payload - send ID in the http path
	v1Router.Delete("/feed_follows/{feedFollowID}", apiConfig.middlewareAuth(apiConfig.handlerDeleteFeedFollows))

	v1Router.Get("/posts", apiConfig.middlewareAuth(apiConfig.handlerGetPostsForUser))

	router.Mount("/v1", v1Router)

	// Connect router to HTTP server
	log.Printf("Server starting on port %v", portStr)
	serv := &http.Server{
		Handler: router,
		Addr:    ":" + portStr,
	}
	err = serv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
