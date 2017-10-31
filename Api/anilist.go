package Api

import (
	"net/http"
	"net/url"
	"log"
	"encoding/json"
)

const (
	ANILIST_BASE_URL = "https://anilist.co/api/"
	ANILIST_AUTH_URL = ANILIST_BASE_URL + "auth/access_token"
)

type Anime struct {

}

type AniListClient struct{
	client_id string
	client_secret string
	access_token string
	expiry_date int
}

func (c *AniListClient) RefreshToken() {

	// Format Query Params
	auth_url, _ := url.Parse(ANILIST_AUTH_URL)
	q := auth_url.Query()
	q.Add("grant_type", "client_credentials")
	q.Add("client_id", c.client_id)
	q.Add("client_secret", c.client_secret)
	auth_url.RawQuery = q.Encode()

	// Request Authentication Token
	req, err := http.Post(auth_url.String(), "application/json", nil)
	defer req.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	// Parse JSON response
	type AniListAuthResponse struct{
		token_type string
		access_token string
		expires_in int
		expires int
	}

	resp := &AniListAuthResponse{}
	err = json.NewDecoder(req.Body).Decode(&resp)
	if err != nil {
		log.Fatal(err)
	}

	c.expiry_date = resp.expires
	c.access_token = resp.access_token
}

// Create New AniListAPI Client from Credentials
func NewAniListClient(client_id string, client_secret string) (*AniListClient) {
	client := &AniListClient{client_id: client_id, client_secret: client_secret}
	client.RefreshToken()
	return client
}
