package delivery

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	"forum/internal/models"
	"forum/internal/service"
)

const (
	authURL      = "https://accounts.google.com/o/oauth2/auth"
	tokenURL     = "https://oauth2.googleapis.com/token"
	clientID     = "817065491333-kpo61kgnd5s5fj2teec254mm48ie70v8.apps.googleusercontent.com"
	clientSecret = "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"
)

var (
	googleSignInConfig = &oauthConfig{
		clientID:     clientID,
		clientSecret: clientSecret,
		scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		redirectURL:  "http://localhost:8080/sign-in/google/callback",
	}
	googleSignUpConfig = &oauthConfig{
		clientID:     clientID,
		clientSecret: clientSecret,
		scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		redirectURL:  "http://localhost:8080/sign-up/google/callback",
	}
)

func requestToGoogle(w http.ResponseWriter, r *http.Request, cfg *oauthConfig) {
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
	requestToGoogle(w, r, googleSignInConfig)
}

func (h *Handler) googleSignUp(w http.ResponseWriter, r *http.Request) {
	requestToGoogle(w, r, googleSignUpConfig)
}

func (h *Handler) signInCallbackFromGoogle(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	user, err := userFromGoogleInfo(code, googleSignInConfig)
	if err != nil {
		h.errorPage(w, http.StatusUnauthorized, err)
		return
	}
	if err := h.setSession(w, user); err != nil {
		if errors.Is(err, service.ErrNoUser) || errors.Is(err, service.ErrWrongPassword) {
			h.errorPage(w, http.StatusUnauthorized, err)
			return
		}
		h.errorPage(w, http.StatusInternalServerError, err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) signUpCallbackFromGoogle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := models.TemplateData{
			Template: "oauth2",
		}
		if err := h.tmpl.ExecuteTemplate(w, "base", data); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}
	case http.MethodPost:
		code := r.URL.Query().Get("code")

		user, err := userFromGoogleInfo(code, googleSignUpConfig)
		if err != nil {
			h.errorPage(w, http.StatusUnauthorized, err)
			return
		}

		if err := r.ParseForm(); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}
		username := r.Form["username"]

		user.Username = username[0]

		if err := h.services.Authorization.CreateUser(*user); err != nil {
			if errors.Is(err, service.ErrEmailTaken) || errors.Is(err, service.ErrUsernameTaken) {
				h.errorPage(w, http.StatusBadRequest, err)
				return
			}
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}

		if err := h.setSession(w, user); err != nil {
			if errors.Is(err, service.ErrNoUser) {
				h.errorPage(w, http.StatusUnauthorized, err)
				return
			}
			h.errorPage(w, http.StatusInternalServerError, err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.errorPage(w, http.StatusMethodNotAllowed, nil)
	}
}

func userFromGoogleInfo(code string, cfg *oauthConfig) (*models.User, error) {
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
		Email:      u.Email,
		AuthMethod: "google",
	}

	return user, nil
}

func googleAccessToken(cfg *oauthConfig, code string) (string, error) {
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
