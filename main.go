package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello world")

	godotenv.Load(".env")
	portStr := os.Getenv("PORT")

	if portStr == "" {
		log.Fatal("PORT is not found in the env")
	}
	fmt.Println("Port: ", portStr)

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

	router.Mount("/v1", v1Router)

	// Connect router to HTTP server
	log.Printf("Server starting on port %v", portStr)
	serv := &http.Server{
		Handler: router,
		Addr:    ":" + portStr,
	}
	err := serv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
