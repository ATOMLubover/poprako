package value

type LoginUserResult struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
}

type RegisterUserResult struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
}
