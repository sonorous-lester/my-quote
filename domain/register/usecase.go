package register

type Usecase interface {
	Register(user NewUser) error
}
