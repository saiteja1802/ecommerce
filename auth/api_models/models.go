package api_models

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r SignupRequest) GetEmail() string {
	return r.Email
}

func (r SignupRequest) GetPassword() string {
	return r.Password
}

type SignupResponse struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func (r SignupResponse) GetUserID() string {
	return r.UserID
}

func (r SignupResponse) GetEmail() string {
	return r.Email
}

func (r SignupResponse) GetCreatedAt() string {
	return r.CreatedAt
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r LoginRequest) GetEmail() string {
	return r.Email
}

func (r LoginRequest) GetPassword() string {
	return r.Password
}

type LoginResponse struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

func (r LoginResponse) GetUserID() string {
	return r.UserID
}

func (r LoginResponse) GetToken() string {
	return r.Token
}

type AuthenticateRequest struct {
	Token string
}

func (r AuthenticateRequest) GetToken() string {
	return r.Token
}

type AuthenticateResponse struct {
	UserID string `json:"user_id"`
}

func (r AuthenticateResponse) GetUserID() string {
	return r.UserID
}
