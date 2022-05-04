package auth

type NewUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Anonymous struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
