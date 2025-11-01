package configuration

import "github.com/daeroworld/shared/configuration"

type Variable struct {
	Database *configuration.Database
	Api      *configuration.Api
	Frontend *configuration.Api
}

func NewVariable() *Variable {
	return &Variable{
		Database: &configuration.Database{
			Uri:      "127.0.0.1",
			Username: "root",
			Password: "root",
		},
		Api: &configuration.Api{
			Port: 25013,
		},
		Frontend: &configuration.Api{
			Ip:   "localhost",
			Port: 25001,
		},
	}
}
