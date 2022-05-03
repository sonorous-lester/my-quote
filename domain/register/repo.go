package register

type Repository interface {
	FindUser(email string) (bool, error)
	Register(name string, email string, password string) error
}
