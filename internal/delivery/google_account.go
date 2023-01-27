package delivery

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	"forum/internal/models"
)

const (
	authURL      = "https://accounts.google.com/o/oauth2/auth"
	tokenURL     = "https://oauth2.googleapis.com/token"
	clientID     = "817065491333-kpo61kgnd5s5fj2teec254mm48ie70v8.apps.googleusercontent.com"
	clientSecret = "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"
)

var (
	googleSignInConfig = &oauthGoogleCfg{
		clientID:     clientID,
		clientSecret: clientSecret,
		scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		redirectURL:  "http://localhost:8080/sign-in/google/callback",
	}
	googleSignUpConfig = &oauthGoogleCfg{
		clientID:     clientID,
		clientSecret: clientSecret,
		scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		redirectURL:  "http://localhost:8080/sign-up/google/callback",
	}
)

func requestToGoogle(w http.ResponseWriter, r *http.Request, cfg oauthGoogleCfg) {
	URL, err := url.Parse(authURL)
	if err != nil {
		log.Printf("Parse: %s", err)
	}

	parameters := url.Values{}
	parameters.Add("client_id", cfg.clientID)
	parameters.Add("redirect_uri", cfg.redirectURL)
	parameters.Add("scope", strings.Join(cfg.scopes, " "))
	parameters.Add("response_type", "code")
	URL.RawQuery = parameters.Encode()
	url := URL.String()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) googleSignIn(w http.ResponseWriter, r *http.Request) {
	requestToGoogle(w, r, *googleSignInConfig)
}

func (h *Handler) googleSignUp(w http.ResponseWriter, r *http.Request) {
	requestToGoogle(w, r, *googleSignUpConfig)
}

func (h *Handler) signInCallbackFromGoogle(w http.ResponseWriter, r *http.Request) {
	user, err := userFromGoogleInfo(r, googleSignInConfig)
	if err != nil {
		h.errorPage(w, http.StatusUnauthorized, err)
		return
	}
	h.setSession(w, user, true)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) signUpCallbackFromGoogle(w http.ResponseWriter, r *http.Request) {
	user, err := userFromGoogleInfo(r, googleSignUpConfig)
	if err != nil {
		h.errorPage(w, http.StatusUnauthorized, err)
		return
	}
	if err := h.services.Authorization.CreateUser(*user, true); err != nil {
		log.Println(err)
		h.errorPage(w, http.StatusUnauthorized, err)
		return
	}

	h.setSession(w, user, true)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func userFromGoogleInfo(r *http.Request, cfg *oauthGoogleCfg) (*models.User, error) {
	code := r.FormValue("code")
	accessToken, err := googleAccessToken(cfg, code)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"GET",
		"https://www.googleapis.com/oauth2/v2/userinfo?access_token="+url.QueryEscape(accessToken),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var u *UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, err
	}

	if u.Email == "" {
		return nil, errors.New("email is empty")
	}

	user := &models.User{
		Username: strings.Split(u.Email, "@")[0],
		Email:    u.Email,
	}

	return user, nil
}

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
