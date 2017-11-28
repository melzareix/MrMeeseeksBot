package Api

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/melzareix/MrMeeseeksBot/Backend/Database"
	"github.com/melzareix/MrMeeseeksBot/Backend/Models"
	"google.golang.org/api/calendar/v3"
	"github.com/grokify/html-strip-tags-go"
)

func ChatHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		err := Models.Error{
			Status:  true,
			Code:    http.StatusOK,
			Message: "Preflight Request!"}
		err.ErrorAsJSON(w)
		return
	}

	if r.Method != http.MethodPost {
		err := Models.Error{
			Status:  false,
			Code:    http.StatusMethodNotAllowed,
			Message: r.Method + " Method Not Allowed. Only POST requests are allowed."}
		err.ErrorAsJSON(w)
		return
	}

	userUuid := r.Header.Get("authorization")
	user, err := Database.DB.GetUser(userUuid)
	if err != nil {
		err := Models.Error{
			Status:  false,
			Code:    http.StatusBadRequest,
			Message: "Invalid user UUID."}
		err.ErrorAsJSON(w)
		return
	}

	var resp map[string]string

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		log.Fatal(err)
	}

	message := resp["message"]
	HandleMessage(message, user, w)
}

func HandleMessage(message string, user *Models.User, w http.ResponseWriter) {
	msg := strings.Split(strings.ToLower(message), " ")
	command := msg[0]

	switch command {
	case "schedule":
		HandleScheduling(strings.Join(msg[1:], " "), user, w)
	case "recommend":
		HandleRecommendation(strings.Join(msg[1:], " "), w)
	case "show":
		HandleAnimeDetails(strings.Join(msg[2:], " "), w)
	default:
		err := Models.Error{
			Status:  false,
			Code:    http.StatusBadRequest,
			Message: "Command not recognized!\n\n" + COMMANDS}
		err.ErrorAsJSON(w)
	}
}

func HandleScheduling(name string, user *Models.User, w http.ResponseWriter) {
	client, err := NewAniListClient("", "")
	if err != nil {
		err := Models.Error{
			Status:  false,
			Code:    http.StatusInternalServerError,
			Message: "Failed to connect to API."}
		err.ErrorAsJSON(w)
		return
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

	// No Token
	// Send Response With OAuth link
	if u.Token == nil {
		config, err := u.GetConfig()
		if err != nil {
			err := Models.Error{
				Status:  false,
				Code:    http.StatusUnauthorized,
				Message: "Failed to Authorize with Google Calendar.",
			}
			err.ErrorAsJSON(w)
		}

		authUrl := u.generateTokenUrl(config)
		resp := Models.Error{
			Status:  false,
			Code:    http.StatusUnauthorized,
			Message: "Oops! You didn't link your Google Calendar Account!\n\n Click this url to link it\n\n" + authUrl,
		}
		resp.ErrorAsJSON(w)
		return
	}

	event := &calendar.Event{
		Summary: results[0].TitleEnglish + " Episode " + strconv.Itoa(selectedIndex),
		Start: &calendar.EventDateTime{
			DateTime: selectedTime.Format(time.RFC3339),
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

	eventLink, err := u.AddEvent(event)
	if err != nil {
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
	formattedTime := humanize.Time(selectedTime)
	eventLink = "<a style='color:black' target='_blank' href='" + eventLink + "'>" + eventLink + "</a>"
	resp.Message = "🕐 Next Episode " + formattedTime + "\n" + strip.StripTags(eventLink)

	RespondWithJSON(w, &resp)
}

func HandleRecommendation(name string, w http.ResponseWriter) {
	client, err := NewAniListClient("", "")

	if err != nil {
		err := Models.Error{
			Status:  false,
			Code:    http.StatusInternalServerError,
			Message: "Failed to connect to API."}
		err.ErrorAsJSON(w)
		return
	}
	results, err := client.Search(name)

	if err != nil {
		err := Models.Error{
			Status:  false,
			Code:    http.StatusBadRequest,
			Message: "No Recommendations for " + name + "."}
		err.ErrorAsJSON(w)
		return
	}

	genreNumber := Randomize(len(results[0].Genres))
	recommendations, _ := client.Recommended(results[0].Genres[genreNumber])

	recommendedNumber := Randomize(len(recommendations))
	recommendedAnime := recommendations[recommendedNumber]

	resp := Models.RecommendationResponse{}
	animeListUrl := "https://myanimelist.net/search/all?q=" + url.QueryEscape(recommendedAnime.TitleEnglish)
	resp.Status = true
	resp.Code = http.StatusOK
	resp.AnimeTitle = recommendedAnime.TitleEnglish
	resp.ImageURL = recommendedAnime.ImageUrlLge
	resp.Message = "More Information:\n" + animeListUrl
	RespondWithJSON(w, &resp)

}

func HandleAnimeDetails(name string, w http.ResponseWriter) {
	client, err := NewAniListClient("", "")

	if err != nil {
		err := Models.Error{
			Status:  false,
			Code:    http.StatusInternalServerError,
			Message: "Failed to connect to API."}
		err.ErrorAsJSON(w)
		return
	}
	results, err := client.Search(name)

	if err != nil {
		err := Models.Error{
			Status:  false,
			Code:    http.StatusBadRequest,
			Message: "No Results for Anime " + name + "."}
		err.ErrorAsJSON(w)
		return
	}

	if len(results) == 0 {
		err := Models.Error{
			Status:  false,
			Code:    http.StatusBadRequest,
			Message: "No Results for Anime " + name + "."}
		err.ErrorAsJSON(w)
		return
	}

	top := results[0]
	resp := Models.AnimeDetailResponse{}
	resp.Status = true
	resp.Code = http.StatusOK
	resp.AnimeTitle = top.TitleEnglish
	resp.Message = fmt.Sprintf(
		"%s\n\nTotal Episodes: %d\n\nDuration: %d\n\nAiring Status: %s\n\nAverage Score: %f",
		strip.StripTags(top.Description), top.TotalEpisodes,
		top.Duration, top.AiringStatus, top.AverageScore)
	resp.ImageURL = top.ImageUrlLge
	RespondWithJSON(w, &resp)
}

func Randomize(upperBound int) (result int) {
	result = int(rand.Float64() * float64(upperBound))
	return
}
