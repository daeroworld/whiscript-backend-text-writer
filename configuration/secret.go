package configuration

import "github.com/daeroworld/shared/configuration"

type Variable struct {
	Database       *configuration.Database
	Api            *configuration.Api
	VoiceReaderApi *configuration.Api
	Frontend       *configuration.Api
}

func NewVariable() *Variable {
	sharedVariable := configuration.NewVariable()
	return &Variable{
		Database: &configuration.Database{
			Uri:      "127.0.0.1",
			Username: "root",
			Password: "root",
		},
		Api: &configuration.Api{
			Port: sharedVariable.TextWriterApi.Port,
		},
		VoiceReaderApi: &configuration.Api{
			Ip:   "{VOICE_READER_IP}",
			Port: 25012,
		},
		Frontend: &configuration.Api{
			Ip:   "localhost",
			Port: sharedVariable.Frontend.Port,
		},
	}
}
