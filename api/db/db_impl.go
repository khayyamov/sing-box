package db

type ImplementationDb struct {
	Interface DbInterface
}

func (pr *ImplementationDb) UserIdIsRequired() string {
	return "Hello"
}
