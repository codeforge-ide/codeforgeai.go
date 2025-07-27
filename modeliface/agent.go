package modeliface

type Agent interface {
	Name() string
	Login() error
	Logout() error
	IsAuthenticated() bool
}
