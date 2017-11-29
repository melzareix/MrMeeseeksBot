package Api

import (
	"net/http"

	"github.com/melzareix/MrMeeseeksBot/Backend/Database"
	"github.com/melzareix/MrMeeseeksBot/Backend/Models"
	//"encoding/json"
)

const (
	WELCOME_MESSAGE = "I'm Mr Meseeks look at me!\n\n I can show you info about an anime, " +
		"recommend and schedule an anime.\n\n" + COMMANDS
	COMMANDS = "Example Commands:\n" +
		"Info Death Note\n" +
		"Recommend Death Note\n" +
		"Schedule Inuyashiki"
)

// Handle the WelcomeResponse Route
func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		err := Models.Error{
			Status:  false,
			Code:    http.StatusMethodNotAllowed,
			Message: r.Method + " Method Not Allowed. Only GET requests are allowed."}
		err.ErrorAsJSON(w)
		return
	}

	user := Models.NewUser()
	err := Database.DB.CreateUser(user)

	if err != nil {
		err := Models.Error{
			Status:  false,
			Code:    http.StatusBadRequest,
			Message: "Failed to create user."}
		err.ErrorAsJSON(w)
		return
	}

	u := Models.WelcomeResponse{Uuid: user.Uuid}
	u.Status = true
	u.Code = http.StatusOK
	u.Message = WELCOME_MESSAGE

	RespondWithJSON(w, &u)
}
