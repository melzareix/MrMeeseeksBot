package Api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
	"log"
	"strconv"
	"os"
)

const (
	ANILIST_BASE_URL   = "https://anilist.co/api/"
	ANILIST_AUTH_URL   = ANILIST_BASE_URL + "auth/access_token"
	ANILIST_SEARCH_URL = ANILIST_BASE_URL + "anime/search/"
	ANILIST_AIRING_URL = ANILIST_BASE_URL + "anime/"
	ANILIST_GENRE_SEARCH = ANILIST_BASE_URL + "browse/anime?genres="
)

type Anime struct {
	Id            int      `json:"id"`
	TitleEnglish string `json:"title_english"`
	Genres        []string `json:"genres"`
	ImageUrlMed   string   `json:"image_url_med"`
	ImageUrlLge   string    `json:"image_url_lge"`
	AiringStatus  string   `json:"airing_status"`
	TotalEpisodes int      `json:"total_episodes"`
	Duration      int64      `json:"duration"`
}

type AniListClient struct {
	client_id     string
	client_secret string
	access_token  string
	expiry_date   int64
}

// Refresh grant token
func (c *AniListClient) RefreshToken() error {

	// Check that current token didn't expire
	if c.expiry_date > 0 {
		expiry_date := time.Unix(c.expiry_date, 0)
		t := time.Now()
		d := expiry_date.Sub(t)
		if d > 0 {
			return nil
		}
	}

	// Format Query Params
	q := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     c.client_id,
		"client_secret": c.client_secret,
	}
	auth_url := FormatUrl(ANILIST_AUTH_URL, q)

	// Request Authentication Token
	req, err := http.Post(auth_url, "application/json", nil)
	defer req.Body.Close()

	if err != nil {
		return err
	}

	// Parse JSON response
	type AniListAuthResponse struct {
		Token_type   string `json:"token_type"`
		Access_token string `json:"access_token"`
		Expires_in   int    `json:"expires_in"`
		Expires      int64  `json:"expires"`
	}

	resp := &AniListAuthResponse{}
	err = json.NewDecoder(req.Body).Decode(&resp)
	if err != nil {
		return err
	}
	c.expiry_date = resp.Expires
	c.access_token = resp.Access_token

	return nil
}

// Search for a certain Anime
func (c *AniListClient) Search(name string) ([]Anime, error) {

	search_url :=  ANILIST_SEARCH_URL + url.PathEscape(name)
	search_url, err := setAccessToken(c, search_url)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(search_url)
	if err != nil {
		return nil, err
	}

	var results []Anime
	log.Println(search_url)

	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (c *AniListClient) Recommended(genre string) ([]Anime, error) {

	url :=  ANILIST_GENRE_SEARCH + url.PathEscape(genre)

	url, err := setAccessToken(c, url)

	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	var results []Anime
	log.Println(url)

	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Get Episodes Airing Dates
func (c *AniListClient) GetAiringDates(id int) (map[string]int64, error) {
	airing_url := ANILIST_AIRING_URL + url.PathEscape(strconv.Itoa(id)) + "/airing"
	airing_url, err := setAccessToken(c, airing_url)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(airing_url)
	if err != nil {
		return nil, err
	}

	var results map[string]int64

	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return nil, err
	}
	return results, nil

}

// Create New AniListAPI Client from Credentials
func NewAniListClient(client_id string, client_secret string) (*AniListClient, error) {
	if client_id == "" && client_secret == "" {
		client_id = os.Getenv("ANILIST_CLIENT_ID")
		client_secret = os.Getenv("ANILIST_CLIENT_SECRET")
	}
	client := &AniListClient{client_id: client_id, client_secret: client_secret}
	err := client.RefreshToken()
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Format URL params
func FormatUrl(base_url string, q_params map[string]string) string {
	u, _ := url.Parse(base_url)
	q := u.Query()

	for k, v := range q_params {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// Set Access Token
func setAccessToken(c *AniListClient, u string) (string, error) {
	err := c.RefreshToken()
	if err != nil {
		return "", err
	}

	new_url, err := url.Parse(u)
	if err != nil {
		return "", nil
	}

	q := new_url.Query()
	q.Add("access_token", c.access_token)

	new_url.RawQuery = q.Encode()
	return new_url.String(), nil
}