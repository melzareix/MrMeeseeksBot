package Api

import (
	"context"
	//"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/melzareix/MrMeeseeksBot/Backend/Database"
	"github.com/melzareix/MrMeeseeksBot/Backend/Models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type CalendarUser struct {
	*Models.User
	srv *calendar.Service
}

// Google Calendar Events
func (u *CalendarUser) generateTokenUrl(config *oauth2.Config) string {
	return config.AuthCodeURL(u.Uuid, oauth2.AccessTypeOffline)
}

// Authorize OAUTH2 Google calendar
func (u *CalendarUser) getToken(config *oauth2.Config, code string) error {
	tok, err := config.Exchange(context.TODO(), code)
	if err != nil {
		return err
	}
	u.Token = tok
	err = Database.DB.SaveUser(u.User)
	return err
}

func (u *CalendarUser) GetConfig() (*oauth2.Config, error) {
	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (u *CalendarUser) SetCalendarService() error {
	ctx := context.Background()
	config, err := u.GetConfig()
	if err != nil {
		return err
	}

	client := config.Client(ctx, u.Token)
	srv, err := calendar.New(client)
	if err != nil {
		return err
	}
	u.srv = srv
	return nil
}

// Add Event to google calendar
func (u *CalendarUser) AddEvent(event *calendar.Event) (string, error) {
	calendarId := "primary"
	event, err := u.srv.Events.Insert(calendarId, event).Do()
	if err != nil {
		return "", err
	}
	return event.HtmlLink, nil
}

// Handle the Calendar Authorization
func CalendarAuthorizationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		err := Models.Error{
			Status:  false,
			Code:    http.StatusMethodNotAllowed,
			Message: r.Method + " Method Not Allowed. Only GET requests are allowed."}
		err.ErrorAsPlainText(w)
		return
	}

	userUuid := r.Header.Get("authorization")
	user, err := Database.DB.GetUser(userUuid)
	if err != nil {
		err := Models.Error{
			Status:  false,
			Code:    http.StatusBadRequest,
			Message: "User " + userUuid + " not found."}
		err.ErrorAsPlainText(w)
		return
	}

	calendarUser := CalendarUser{User: user}
	config, err := calendarUser.GetConfig()
	if err != nil {
		err := Models.Error{
			Status:  false,
			Code:    http.StatusBadRequest,
			Message: err.Error()}
		err.ErrorAsPlainText(w)
		return
	}
	http.Redirect(w, r, calendarUser.generateTokenUrl(config), http.StatusTemporaryRedirect)
}

// Callback For Authorization
func CalendarAuthorizationCallbackHandler(w http.ResponseWriter, r *http.Request) {
	userUuid := r.FormValue("state")
	code := r.FormValue("code")

	user, err := Database.DB.GetUser(userUuid)
	if err != nil {
		err := Models.Error{
			Status:  false,
			Code:    http.StatusBadRequest,
			Message: "User " + userUuid + " not found."}
		err.ErrorAsPlainText(w)
		return
	}

	u := CalendarUser{User: user}
	config, err := u.GetConfig()
	if err != nil {
		err := Models.Error{
			Status:  false,
			Code:    http.StatusBadRequest,
			Message: err.Error()}
		err.ErrorAsPlainText(w)
		return
	}

	err = u.getToken(config, code)
	if err != nil {
		err := Models.Error{
			Status:  false,
			Code:    http.StatusBadRequest,
			Message: err.Error()}
		err.ErrorAsPlainText(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "<html><body><strong>Google Calendar Authorized! You may close this window "+
		"and go back to the chat bot!</body></html>")
}
