package master

type Storage struct {
	Console
}

func New() *Storage {
	return &Storage{Console{}}
}
