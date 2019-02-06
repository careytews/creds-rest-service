package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const tokeninfoEndpointURL = "https://www.googleapis.com/oauth2/v3/tokeninfo?access_token="
const clientID = "1041863416400-luuj7j2h8a1mdi454hf2lnqqngfrbevc.apps.googleusercontent.com"

// AccessTokenInfo information about an access token
type AccessTokenInfo struct {

	// The client_id of the authorized presenter.
	AuthorizedPresenter string `json:"azp,omitempty"`

	// Identifies the audience that this ID token is intended for.
	Audience string `json:"aud,omitempty"`

	// The subject of the token.An identifier for the user, unique among all Google accounts and never reused.
	Subject string `json:"sub,omitempty"`

	// The access scope of the token
	Scope string `json:"scope,omitempty"`

	// The time the ID token expires, represented in Unix time
	ExpireTime string `json:"exp,omitempty"`

	// The lifetime of the token, in seconds
	ExpiresIn string `json:"expires_in,omitempty"`

	// The user's email address
	Email string `json:"email,omitempty"`

	// "true" if the user's e-mail address has been verified; otherwise "false".
	EmailVerified string `json:"email_verified,omitempty"`

	// Indicates whether the application can refresh access tokens when the user is not present at the browser.
	AccessType string `json:"access_type,omitempty"`

	// Indicates an error has occurred if not empty
	ErrorDescription string `json:"error_description,omitempty"`
}

// GetTokenInfo gets info about the token
func getTokenInfo(token string) (*AccessTokenInfo, error) {

	// Set a timeout in case of errors that block for ages
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(tokeninfoEndpointURL + token)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tokenInfo := AccessTokenInfo{}

	err = json.Unmarshal(body, &tokenInfo)
	if err != nil {
		return nil, err
	}

	err = validateAccessToken(&tokenInfo)

	return &tokenInfo, err
}

// validateAccessToken checks an AccessTokenInfo and responds if there's an error
func validateAccessToken(tokenInfo *AccessTokenInfo) error {

	var err error

	// Check for nil, error description, client ID/audience and then the scopes
	// We only check email if all of the rest have passed, and we will probably
	// not have to check the email very often
	if tokenInfo == nil {
		err = errors.New("tokenInfo is nil")
	} else if len(tokenInfo.ErrorDescription) > 0 {
		err = errors.New(tokenInfo.ErrorDescription)
	} else if tokenInfo.Audience != clientID {
		err = errors.New("wrong client ID")
	} else {
		err = accessTokenHasScope(tokenInfo.Scope)
	}

	if err == nil {
		if tokenInfo.EmailVerified != "true" {
			err = errors.New("email verification failed")
		} else if len(tokenInfo.Email) == 0 {
			err = errors.New("no email address returned")
		}
	}

	return err
}

func accessTokenHasScope(scope string) error {

	// Can't do a const, but you can do a fixed-size array
	requires := [...]string{"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/plus.me"}

	has := true
	var err error
	missing := make([]string, 0)

	for _, required := range requires {
		has = strings.Contains(scope, required)
		if has != true {
			missing = append(missing, required)
		}
	}
	if len(missing) > 0 {
		err = fmt.Errorf("missing scope(s) %s", strings.Join(missing, " "))
	}

	return err
}
