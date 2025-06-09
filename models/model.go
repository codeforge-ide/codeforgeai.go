package models

type Model interface {
	SendRequest(prompt string, config interface{}) (string, error)
}
