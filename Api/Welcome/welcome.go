package Welcome

import (
	"net/http"
	"github.com/melzareix/MrMeseeksBot/Models"
	"github.com/melzareix/MrMeseeksBot/Database"
	"github.com/satori/go.uuid"
	"encoding/json"
)

// Handle the Welcome Route
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		err := Models.Error{
			Status: false,
			Code:    http.StatusMethodNotAllowed,
			Message: r.Method + " Method Not Allowed. Only GET requests are allowed."}
		err.ErrorAsJSON(w)
		return
	}

	uniqueId := uuid.NewV4().String()
	user := Models.User{Uuid: uniqueId, AccessToken: "", TokenType: "", RefreshToken: "", Expiry: ""}
	err := Database.DB.CreateUser(&user)

	if err != nil {
		err := Models.Error{
			Status: false,
			Code:    http.StatusBadRequest,
			Message: "Failed to create user."}
		err.ErrorAsJSON(w)
		return
	}

	u := Models.Welcome{Message: "Welcome!", Uuid: uniqueId}
	RespondWithJSON(w, &u)
}


func RespondWithJSON(w http.ResponseWriter, u *Models.Welcome) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&u)
}