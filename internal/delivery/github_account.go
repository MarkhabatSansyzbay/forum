package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"forum/internal/models"
	"forum/internal/service"
)

const (
	ghAuthURL      = "https://github.com/login/oauth/authorize"
	ghTokenURL     = "https://github.com/login/oauth/access_token"
	ghClientID     = "bf208a1ead8eda54ce05"
	ghClientSecret = "e28468b0b7064f3c4aa44c3660f93bbac4969588"
)

var (
	ghSignInCfg = &oauthConfig{
		clientID:     ghClientID,
		clientSecret: ghClientSecret,
		redirectURL:  "http://localhost:8080/github/callback/sign-in",
		scopes:       []string{"user:email"},
	}
	ghSignUpCfg = &oauthConfig{
		clientID:     ghClientID,
		clientSecret: ghClientSecret,
		redirectURL:  "http://localhost:8080/github/callback/sign-up",
		scopes:       []string{"user:email"},
	}
)

func requestToGithub(w http.ResponseWriter, r *http.Request, cfg oauthConfig) {
	URL, err := url.Parse(ghAuthURL)
	if err != nil {
		log.Printf("Parse: %s", err)
	}

	parameters := url.Values{}
	parameters.Add("client_id", cfg.clientID)
	parameters.Add("redirect_uri", cfg.redirectURL)
	parameters.Add("scope", strings.Join(cfg.scopes, " "))

	URL.RawQuery = parameters.Encode()
	url := URL.String()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) githubSignIn(w http.ResponseWriter, r *http.Request) {
	requestToGithub(w, r, *ghSignInCfg)
}

func (h *Handler) githubSignUp(w http.ResponseWriter, r *http.Request) {
	requestToGithub(w, r, *ghSignUpCfg)
}

func (h *Handler) signInCallbackGithub(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	user, err := h.userFromGithubInfo(code, ghSignInCfg)
	if err != nil {
		h.errorPage(w, http.StatusUnauthorized, err)
		return
	}

	if err := h.setSession(w, user, true); err != nil {
		if errors.Is(err, service.ErrNoUser) || errors.Is(err, service.ErrWrongPassword) {
			h.errorPage(w, http.StatusUnauthorized, err)
			return
		}
		h.errorPage(w, http.StatusInternalServerError, err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) signUpCallbackGithub(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	user, err := h.userFromGithubInfo(code, ghSignUpCfg)
	if err != nil {
		h.errorPage(w, http.StatusUnauthorized, err)
		return
	}

	if err := h.services.Authorization.CreateUser(*user, true); err != nil {
		if errors.Is(err, service.ErrEmailTaken) || errors.Is(err, service.ErrUsernameTaken) {
			h.errorPage(w, http.StatusBadRequest, err)
			return
		}
		h.errorPage(w, http.StatusInternalServerError, err)
		return
	}

	if err := h.setSession(w, user, true); err != nil {
		if errors.Is(err, service.ErrNoUser) {
			h.errorPage(w, http.StatusUnauthorized, err)
			return
		}
		h.errorPage(w, http.StatusInternalServerError, err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) userFromGithubInfo(code string, cfg *oauthConfig) (*models.User, error) {
	accessToken, err := githubAccessToken(cfg, code)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		return nil, err
	}

	authHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authHeaderValue)
	req.Header.Set("accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var u *GithubUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, err
	}

	email, err := emailFromGithub(code, accessToken, cfg)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:   u.Username,
		Email:      email,
		AuthMethod: "github",
	}

	return user, nil
}

func emailFromGithub(code, accessToken string, cfg *oauthConfig) (string, error) {
	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user/emails",
		nil,
	)
	if err != nil {
		return "", err
	}
	authHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authHeaderValue)
	req.Header.Set("accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data []*GithubUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if data[0].Email == "" {
		return "", errors.New("email is empty")
	}

	return data[0].Email, nil
}

func githubAccessToken(cfg *oauthConfig, code string) (string, error) {
	v := url.Values{
		"code":          {code},
		"client_id":     {cfg.clientID},
		"client_secret": {cfg.clientSecret},
	}

	req, err := http.NewRequest("POST", ghTokenURL, strings.NewReader(v.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

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
