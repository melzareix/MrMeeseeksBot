package Api

import (
	"fmt"
	"github.com/melzareix/MrMeeseeksBot/Database"
	"github.com/melzareix/MrMeeseeksBot/Models"
	"google.golang.org/api/calendar/v3"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := Models.Error{
			Status:  false,
			Code:    http.StatusMethodNotAllowed,
			Message: r.Method + " Method Not Allowed. Only POST requests are allowed."}
		err.ErrorAsJSON(w)
		return
	}

	userUuid := r.Header.Get("Authorization")
	user, err := Database.DB.GetUser(userUuid)
	if err != nil {
		err := Models.Error{
			Status:  false,
			Code:    http.StatusBadRequest,
			Message: "Invalid user UUID."}
		err.ErrorAsJSON(w)
		return
	}

	message := r.FormValue("message")
	HandleMessage(message, user, w)
}

func HandleMessage(message string, user *Models.User, w http.ResponseWriter) {
	msg := strings.Split(strings.ToLower(message), " ")
	command := msg[0]

	switch command {
	case "schedule":
		HandleScheduling(strings.Join(msg[1:], " "), user, w)
	}
}

func HandleScheduling(name string, user *Models.User, w http.ResponseWriter) {
	client, err := NewAniListClient("", "")
	if err != nil {
		log.Fatal(err)
	}

	results, err := client.Search(name)
	if err != nil {
		err := Models.Error{
			Status:  false,
			Code:    http.StatusBadRequest,
			Message: "No Results for " + name + "."}
		err.ErrorAsJSON(w)
		return
	}

	// We assume first result is the correct anime
	currentAnime := results[0]
	if currentAnime.AiringStatus != "currently airing" {
		err := Models.Error{
			Status:  false,
			Code:    http.StatusBadRequest,
			Message: "The anime " + name + " has " + currentAnime.AiringStatus + "."}
		err.ErrorAsJSON(w)
		return
	}

	id := currentAnime.Id

	airingDates, err := client.GetAiringDates(id)
	if err != nil {
		err := Models.Error{
			Status:  false,
			Code:    http.StatusBadRequest,
			Message: "Airing Dates not available for " + name + "."}
		err.ErrorAsJSON(w)
		return
	}

	sortedAiringDates := make([]time.Time, len(airingDates))

	for i := 0; i < len(airingDates); i++ {
		sortedAiringDates[i] = time.Unix(airingDates[strconv.Itoa(i+1)], 0)
	}

	allEpisodesFinishedAiring := true
	var selectedIndex int
	var selectedTime time.Time

	for k, v := range sortedAiringDates {
		if v.After(time.Now()) {
			selectedTime = v
			selectedIndex = k + 1
			allEpisodesFinishedAiring = false
			break
		}
	}

	if allEpisodesFinishedAiring {
		err := Models.Error{
			Status:  false,
			Code:    http.StatusBadRequest,
			Message: "All episodes for " + name + "have finished airing."}
		err.ErrorAsJSON(w)
		return
	}

	u := CalendarUser{User: user}
	formattedTime := selectedTime.Format(time.RFC3339)
	event := &calendar.Event{
		Summary: results[0].TitleEnglish + " Episode " + strconv.Itoa(selectedIndex),
		Start: &calendar.EventDateTime{
			DateTime: formattedTime,
		},
		End: &calendar.EventDateTime{
			DateTime: selectedTime.Add(time.Duration(currentAnime.Duration) * time.Minute).Format(time.RFC3339),
		},
	}
	err = u.SetCalendarService()

	if err != nil {
		err := Models.Error{
			Status:  false,
			Code:    http.StatusBadRequest,
			Message: "Failed to connect to google calendar."}
		err.ErrorAsJSON(w)
		return
	}

	_, err = u.AddEvent(event)
	if err != nil {
		fmt.Println(err)
		err := Models.Error{
			Status:  false,
			Code:    http.StatusBadRequest,
			Message: "Failed to add event to google calendar."}
		err.ErrorAsJSON(w)
		return
	}

	resp := Models.SchedulingResponse{}
	resp.Status = true
	resp.Code = http.StatusOK
	resp.Message = "Next Episode at " + formattedTime + " Added to Calendar!\n"

	RespondWithJSON(w, &resp)
}
