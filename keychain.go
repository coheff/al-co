package alco

import (
	"encoding/json"

	"github.com/deanishe/awgo/keychain"
	"golang.org/x/oauth2"
)

const key = "token"

// CacheToken adds a token to given Keychain. If a token already exists, it is replaced.
func CacheToken(kc *keychain.Keychain, tok *oauth2.Token) error {
	jToken, err := json.Marshal(tok)
	if err != nil {
		return err
	}

	err = kc.Set(key, string(jToken))
	if err != nil {
		return err
	}
	return nil
}

// CachedToken retrieves a token from a given Keychain.
func CachedToken(kc *keychain.Keychain) (*oauth2.Token, error) {
	jToken, err := kc.Get(key)
	if err != nil {
		return nil, err
	}

	var tok oauth2.Token
	err = json.Unmarshal([]byte(jToken), &tok)
	if err != nil {
		return nil, err
	}
	return &tok, err
}
