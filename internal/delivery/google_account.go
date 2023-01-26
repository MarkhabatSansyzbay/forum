package delivery

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	authURL      = "https://accounts.google.com/o/oauth2/auth"
	tokenURL     = "https://oauth2.googleapis.com/token"
	clientID     = "817065491333-kpo61kgnd5s5fj2teec254mm48ie70v8.apps.googleusercontent.com"
	clientSecret = "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"
)

type oauthGoogleCfg struct {
	clientID     string
	clientSecret string
	redirectURL  string
	scopes       []string
}

var googleConfig = &oauthGoogleCfg{
	clientID:     clientID,
	clientSecret: clientSecret,
	scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	redirectURL:  "http://localhost:8080/sign-in/google/callback",
}

func (h *Handler) googleSignIn(w http.ResponseWriter, r *http.Request) {
	URL, err := url.Parse(authURL)
	if err != nil {
		log.Printf("Parse: %s", err)
	}

	parameters := url.Values{}
	parameters.Add("client_id", googleConfig.clientID)
	parameters.Add("redirect_uri", googleConfig.redirectURL)
	parameters.Add("scope", strings.Join(googleConfig.scopes, " "))
	parameters.Add("response_type", "code")
	URL.RawQuery = parameters.Encode()
	url := URL.String()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type,omitempty"`
}

func (h *Handler) callbackFromGoogle(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	if code == "" {
		// ?
		w.Write([]byte("Code Not Found to provide AccessToken..\n"))
		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			w.Write([]byte("User has denied Permission.."))
		}
	} else {
		accessToken, err := googleAccessToken(googleConfig, code)
		if err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}

		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(accessToken))
		if err != nil {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		body1, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			// do smth
		}

		fmt.Println(string(body1))
	}
}

// func googleUserInfo() ()

func googleAccessToken(cfg *oauthGoogleCfg, code string) (string, error) {
	v := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {cfg.redirectURL},
		"client_id":     {cfg.clientID},
		"client_secret": {cfg.clientSecret},
	}

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(v.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var token Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", err
	}

	return token.AccessToken, nil
}
