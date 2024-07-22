package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Server responded with 5XX error: ", msg)
	}
	type errResponse struct {
		Error string `json:"error"` // key for the field is "error"
	}
	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, error := json.Marshal(payload) // attempt to marshal payload into JSON string
	if error != nil {
		log.Printf("Failed to marshal JSON: %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json") // add content type response header
	w.WriteHeader(code)
	w.Write(data) // write response body
}
