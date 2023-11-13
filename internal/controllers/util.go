package controllers

import (
	"encoding/json"
	"net/http"

)

func writeJSONResponse(w http.ResponseWriter, statusCode int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

func handleError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	responseJSON, _ := json.Marshal(map[string]string{"error": message})
	w.Write(responseJSON)
}

