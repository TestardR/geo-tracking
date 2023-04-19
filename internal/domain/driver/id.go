package driver

type Id struct {
	id string
}

func NewId(id string) Id {
	return Id{id: id}
}
