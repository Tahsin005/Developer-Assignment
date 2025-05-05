package utils

import (
	"encoding/json"
	"net/http"
)

func WriteError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}