package alco

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os/exec"

	"github.com/deanishe/awgo/keychain"
	"golang.org/x/oauth2"
)

// Token attempts to retrieve a cached OAuth2 token for a given keychain.
// If none exists the OAuth2 flow is started for a given oauth2.Config.
// If a token is successfully retrieved it is cached.
func Token(config *oauth2.Config, kc *keychain.Keychain) *oauth2.Token {
	tok, err := cachedToken(kc)
	if err != nil {
		log.Printf("Error retrieving cached token; it might not exist: %v", err)

		// get new token
		tok, err = newToken(config)
		if err != nil {
			log.Fatalf("Error aquiring token: %v", err)
		}

		// store token
		err = cacheToken(kc, tok)
		if err != nil {
			log.Fatalf("Error storing token: %v", err)
		}
	}

	return tok
}

// newToken starts the OAuth2 flow in order to get a token.
func newToken(config *oauth2.Config) (*oauth2.Token, error) {
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
	_, err := url.ParseRequestURI(redirectUrl)
	if err != nil {
		return nil, err
	}

	listener, err := net.Listen("tcp", redirectUrl)
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
