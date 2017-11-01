package Api

import (
	"encoding/json"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, u interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&u)
}