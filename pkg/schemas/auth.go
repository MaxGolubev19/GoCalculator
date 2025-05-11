package schemas

type contextKey string

const UserContextKey = contextKey("user")

type User struct {
	Login          string
	HashedPassword string
}

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
