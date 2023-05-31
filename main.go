package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	ory "github.com/ory/client-go"
	"golang.org/x/oauth2"
)

var (
	config oauth2.Config
	state  = "randomstate"
	client = "testclient"
)

// Use this context to access Ory APIs which require an Ory API Key.
var oryAuthedContext = context.WithValue(context.Background(), ory.ContextAccessToken, "ory_pat_tXlKB9RHSxZJaBpidIXEdayqpR3RzRQK")

func handleMain(w http.ResponseWriter, r *http.Request) {
	url := config.AuthCodeURL(state)
	fmt.Println("Redirecting user to authorization endpoint:", url)
	http.Redirect(w, r, url, http.StatusFound)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != state {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}

	fmt.Println("Received callback from authorization endpoint with state:", r.URL.Query().Get("state"))
	fmt.Println("Exchanging authorization code for access token")

	token, err := config.Exchange(context.Background(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Received access token and refresh token")
	fmt.Fprintf(w, "Access Token: %s\n", token.AccessToken)
	fmt.Fprintf(w, "Refresh Token: %s\n", token.RefreshToken)
	fmt.Fprintf(w, "Expiry: %s\n", token.Expiry)
}

func main() {
	fmt.Println("-------------------")
	fmt.Println("OAuth2 client is being created")
	fmt.Println("-------------------")

	// Create a new OAuth2 client
	// clientName := "grafana_clientss"
	oAuth2Client := *ory.NewOAuth2Client() // OAuth2Client |
	oAuth2Client.SetClientName(client)
	oAuth2Client.SetScope("openid")
	oAuth2Client.SetGrantTypes([]string{"authorization_code", "refresh_token", "client_credentials"})
	oAuth2Client.SetResponseTypes([]string{"code", "id_token", "token"})
	oAuth2Client.SetRedirectUris([]string{"http://localhost:3001/login/generic_oauth"})
	oAuth2Client.SetTokenEndpointAuthMethod("client_secret_post")

	configuration := ory.NewConfiguration()
	configuration.Servers = []ory.ServerConfiguration{
		{
			URL: "

			", // Public API URL
		},
	}
	ory := ory.NewAPIClient(configuration)
	resp, r, err := ory.OAuth2Api.CreateOAuth2Client(oryAuthedContext).OAuth2Client(oAuth2Client).Execute()
	if err != nil {
		switch r.StatusCode {
		case http.StatusConflict:
			fmt.Fprintf(os.Stderr, "Conflict when creating oAuth2Client: %v\n", err)
		default:
			fmt.Fprintf(os.Stderr, "Error when calling `OAuth2Api.CreateOAuth2Client`: %v\n", err)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		}
	}

	fmt.Println("-------------------")
	fmt.Println("OAuth2 client has been created")
	fmt.Println("-------------------")

	config = oauth2.Config{
		ClientID:     resp.GetClientId(),
		ClientSecret: resp.GetClientSecret(),
		Scopes:       []string{"openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:4000/oauth2/auth",
			TokenURL: "http://localhost:4000/oauth2/token",
		},
		RedirectURL: resp.GetRedirectUris()[0],
	}

	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login/generic_oauth", handleCallback)

	log.Println("Server started on :3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}
