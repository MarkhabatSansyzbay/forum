package delivery

type oauthConfig struct {
	clientID     string
	clientSecret string
	redirectURL  string
	scopes       []string
}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type,omitempty"`
	Scope       string `json:"scope"`
}

type UserInfo struct {
	Email string `json:"email"`
}
