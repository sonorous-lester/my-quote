package register

type Repository interface {
	FindUser(email string) (bool, error)
}
