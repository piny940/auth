package di

type registry struct {
}

var reg = &registry{}

func GetRegistry() *registry {
	return reg
}
