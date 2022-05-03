package common

type Validator interface {
	Validate(s string) bool
}

type HashValidator interface {
	Hash(s string) (string, error)
	Compare(s string, h string) bool
}
