package auth

type Usecase interface {
	Register(user NewUser) error
}
