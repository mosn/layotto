package hello

const ServiceName = "hello"

type HelloService interface {
	Init(*HelloConfig) error
	Hello(*HelloRequest) (*HelloReponse, error)
}

type HelloConfig struct {
	HelloString string `json:"hello"`
}

type HelloRequest struct {
	Name string `json:"name"`
}

type HelloReponse struct {
	HelloString string `json:"hello"`
}
