package handlers

import (
	"encoding/json"
	"net/http"
)

func ResendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Check your e-mail for the verification link",
	})
}