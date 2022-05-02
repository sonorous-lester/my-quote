package register

type Repository interface {
	Find(email string) bool
}
