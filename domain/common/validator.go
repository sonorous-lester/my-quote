package common

type Validator interface {
	Validate(s string) bool
}
