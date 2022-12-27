package alco

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os/exec"

	"golang.org/x/oauth2"
)

// NewToken starts the OAuth2 flow in order to get a token.
func NewToken(config *oauth2.Config) (*oauth2.Token, error) {
	codeCh, err := startWebServer(config.RedirectURL)
	if err != nil {
		return nil, err
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	if exec.Command("open", authURL).Start() != nil {
		return nil, err
	}

	// Wait for the web server to get the code.
	code := <-codeCh

	return exchangeToken(config, code)
}

// startWebServer listens for OAuth2 code returned as part of the three-legged auth flow.
func startWebServer(redirectUrl string) (chan string, error) {
	loopback, err := url.Parse(redirectUrl)
	if err != nil {
		return nil, err
	}

	listener, err := net.Listen("tcp", loopback.Host)
	if err != nil {
		return nil, err
	}

	codeCh := make(chan string)
	go http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		codeCh <- code // send code to OAuth flow
		listener.Close()
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Received code: %v\nYou can now safely close this browser window and start using the workflow.", code)
	}))

	return codeCh, nil
}

// exchangeToken swaps the authorization code for an access token.
func exchangeToken(config *oauth2.Config, code string) (*oauth2.Token, error) {
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	return token, nil
}
