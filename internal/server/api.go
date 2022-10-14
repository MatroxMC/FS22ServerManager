package server

type Api struct {
	Password string
	Username string
	Url      string
}

func NewApi(g Server) (Api, error) {
	d := Api{
		Password: "secret",
		Username: "test",
		Url:      "test",
	}

	return d, nil
}

func (a Api) Status() error {

	return nil
}
