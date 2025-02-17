package core

type ResponseParser interface {
	Parse(*Response) (string, error)
}
