package Models

import (
	"encoding/json"
	"net/http"
)

// JSON Error Message Model

type Error struct {
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Write JSON Error Message
func (err *Error) ErrorAsJSON(w http.ResponseWriter) {
	w.WriteHeader(err.Code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&err)
}