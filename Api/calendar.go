package Api

import (
	"io/ioutil"
	"log"
	"google.golang.org/api/calendar/v3"
	"fmt"
	"golang.org/x/oauth2/google"
	"github.com/melzareix/MrMeeseeksBot/Models"
	"golang.org/x/oauth2"
	"context"
	"net/http"
	"github.com/melzareix/MrMeeseeksBot/Database"
)

type CalendarUser struct {
	*Models.User
}

var (
	srv *calendar.Service
)

// Google Calendar Events
func (u* CalendarUser) generateTokenUrl(config *oauth2.Config) string {
	return config.AuthCodeURL("state-token", oauth2.AccessTypeOffline);
}

func (u* CalendarUser) getToken(config *oauth2.Config, code string) {
	tok, err := config.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	u.Token = tok
	Database.DB.SaveUser(u.User)
}

func (u *CalendarUser) getClient(ctx context.Context, config *oauth2.Config) *http.Client{
	fmt.Println(u.generateTokenUrl(config))
	var code string;
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}
	u.getToken(config, code)
	return config.Client(ctx, u.Token)
}

func (u *CalendarUser) Authorize() {
	ctx := context.Background()
	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Aborting! Cannot read secret file! %v", err)
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
	}

	client := u.getClient(ctx, config)
	srv, err = calendar.New(client)
	if err != nil {
		log.Fatal(err)
	}
}

func (u *CalendarUser) AddEvent(event *calendar.Event) {
	calendarId := "primary"
	event, err := srv.Events.Insert(calendarId, event).Do()
	if err != nil {
		log.Fatalf("Unable to create event. %v\n", err)
	}
	fmt.Printf("Event created: %s\n", event.HtmlLink)
}
