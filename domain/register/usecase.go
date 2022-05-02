package register

type Register interface {
	Process(user NewUser) error
}
