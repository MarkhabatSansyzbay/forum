package delivery

import (
	"fmt"
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
}

var googleConfig = &oauthGoogleCfg{
	clientID:     clientID,
	clientSecret: clientSecret,
	redirectURL:  "http://localhost/sign-in/google/callback",
}

func (h *Handler) googleSignIn(w http.ResponseWriter, r *http.Request) {
	URL, err := url.Parse(authURL)
	if err != nil {
		log.Printf("Parse: %s", err)
	}

	parameters := url.Values{}
	parameters.Add("client_id", googleConfig.clientID)
	parameters.Add("redirect_uri", googleConfig.redirectURL)
	parameters.Add("response_type", "code")
	URL.RawQuery = parameters.Encode()
	url := URL.String()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
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
		v := url.Values{
			"grant_type":    {"authorization_code"},
			"code":          {code},
			"redirect_uri":  {googleConfig.redirectURL},
			"client_id":     {googleConfig.clientID},
			"client_secret": {googleConfig.clientSecret},
		}

		req, err := http.NewRequest("POST", tokenURL, strings.NewReader(v.Encode()))
		if err != nil {
			// do smth
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			// do smth
		}
		fmt.Println(res.Body)
	}
}
