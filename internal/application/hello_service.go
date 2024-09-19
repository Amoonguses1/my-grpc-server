package application

type HelloService struct {
}

func (hs *HelloService) GenerateHello(name string) string {
	return "Hello " + name
}
