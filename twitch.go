package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
)

var (
	clientID = os.Getenv("TWITCH_CLIENT_ID")
	// Consider storing the secret in an environment variable or a dedicated storage system.
	clientSecret = os.Getenv("TWITCH_CLIENT_SECRET")
	oauth2Config *clientcredentials.Config
	httpClient   *http.Client
	pick         = 3
)

func getFollowerCnt() int {
	flag.Parse()

	oauth2Config = &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     twitch.Endpoint.TokenURL,
		Scopes:       []string{"user:read:email"},
	}

	httpClient = oauth2Config.Client(context.TODO())

	// twitch api reference  https://dev.twitch.tv/docs/api/reference

	// 나를 팔로하는 사람들 얻어오기
	u := getUserInfo([]string{"suapapa"}, nil)
	f := getUserFollowTo(u.Data[0].ID)

	return f.Total
}

// User reperesents users info
type User struct {
	Data []struct {
		ID              string `json:"id"`
		Login           string `json:"login"`
		DisplayName     string `json:"display_name"`
		Type            string `json:"type"`
		BroadcasterType string `json:"broadcaster_type"`
		Description     string `json:"description"`
		ProfileImageURL string `json:"profile_image_url"`
		OfflineImageURL string `json:"offline_image_url"`
		ViewCount       int    `json:"view_count"`
	} `json:"data"`
}

func getUserInfo(login []string, id []string) *User {
	values := make(url.Values)
	for _, l := range login {
		values.Add("login", l)
	}
	for _, i := range id {
		values.Add("id", i)
	}
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users?"+values.Encode(), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Client-Id", clientID)

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var user User
	json.NewDecoder(resp.Body).Decode(&user)
	return &user
}

// Follow reperesents users follows
type Follow struct {
	Total int `json:"total"`
	Data  []struct {
		FromID     string    `json:"from_id"`
		FromName   string    `json:"from_name"`
		ToID       string    `json:"to_id"`
		ToName     string    `json:"to_name"`
		FollowedAt time.Time `json:"followed_at"`
	} `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}

func getUserFollowTo(id string) *Follow {
	values := make(url.Values)
	values.Add("to_id", id)
	values.Add("first", "100")
	// values.Add("after", "TODO")
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users/follows?"+values.Encode(), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Client-Id", clientID)

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var follow Follow
	json.NewDecoder(resp.Body).Decode(&follow)
	return &follow
}
