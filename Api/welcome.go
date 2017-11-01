package Api

import (
	"net/http"
	"github.com/melzareix/MrMeeseeksBot/Models"
	"github.com/melzareix/MrMeeseeksBot/Database"
	//"encoding/json"
)

const (
	WELCOME_MESSAGE = "Hello, I'm <b>Mr Meseeks</b> look at me! <br> You can ask me to give you information about " +
		"an anime, Recommend Anime or schedule the next episode of an anime." +
		"<br>==================================<br><b>COMMANDS</b><br>==================================<br>" +
		" 1. Show Information [ANIME NAME]<br>" +
		" 2. Recommend [ANIME NAME]<br>" +
		" 3. Schedule [ANIME NAME]<br>"
)
// Handle the WelcomeResponse Route
func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		err := Models.Error{
			Status: false,
			Code:    http.StatusMethodNotAllowed,
			Message: r.Method + " Method Not Allowed. Only GET requests are allowed."}
		err.ErrorAsPlainText(w)
		return
	}

	user := Models.NewUser()
	err := Database.DB.CreateUser(user)

	if err != nil {
		err := Models.Error{
			Status: false,
			Code:    http.StatusBadRequest,
			Message: "Failed to create user."}
		err.ErrorAsPlainText(w)
		return
	}

	u := Models.WelcomeResponse{Uuid: user.Uuid}
	u.Status = true
	u.Code = http.StatusOK
	u.Message = WELCOME_MESSAGE

	RespondWithJSON(w, &u)
}