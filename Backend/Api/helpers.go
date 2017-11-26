package Api

import (
	"encoding/json"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, u interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&u)
}