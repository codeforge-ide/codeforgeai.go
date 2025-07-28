package modeliface

// BaseAgent provides default implementations for Agent interface methods.
type BaseAgent struct {
	NameStr       string
	Authenticated bool
}

func (b *BaseAgent) Name() string {
	return b.NameStr
}

func (b *BaseAgent) Login() error {
	b.Authenticated = true
	return nil
}

func (b *BaseAgent) Logout() error {
	b.Authenticated = false
	return nil
}

func (b *BaseAgent) IsAuthenticated() bool {
	return b.Authenticated
}
