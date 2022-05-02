package register

type Repository interface {
	FindUser(email string) (bool, error)
	Register(email string, password string) error
}
