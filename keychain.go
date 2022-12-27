package alco

import (
	"encoding/json"

	"github.com/deanishe/awgo/keychain"
	"golang.org/x/oauth2"
)

const key = "token"

// cacheToken adds a token to given Keychain. If a token already exists, it is replaced.
func cacheToken(kc *keychain.Keychain, tok *oauth2.Token) error {
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

// cachedToken retrieves a token from a given Keychain.
func cachedToken(kc *keychain.Keychain) (*oauth2.Token, error) {
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
