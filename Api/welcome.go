package Api

import (
	"net/http"
	"github.com/melzareix/MrMeeseeksBot/Models"
	"github.com/melzareix/MrMeeseeksBot/Database"
	//"encoding/json"
)

const (
	WELCOME_MESSAGE = "Hello, I'm Mr Meseeks look at me ( ͡° ͜ʖ ͡°)!\n You can ask me to give you information about " +
		"an anime, Recommend Anime or schedule the next episode of an anime." +
		"\n==================================\nCOMMANDS\n==================================\n" +
		" 1. Information [ANIME NAME]" +
		" 2. Recommend [ANIME NAME]" +
		" 3. Schedule [ANIME NAME]"
)
// Handle the WelcomeResponse Route
func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		err := Models.Error{
			Status: false,
			Code:    http.StatusMethodNotAllowed,
			Message: r.Method + " Method Not Allowed. Only GET requests are allowed."}
		err.ErrorAsJSON(w)
		return
	}

	user := Models.NewUser()
	err := Database.DB.CreateUser(user)

	if err != nil {
		err := Models.Error{
			Status: false,
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