package main

import (
	"fmt"
	"text/writer/configuration"
)

func main() {
	startContainer()
}

func startContainer() error {
	ctnr := configuration.NewContainer()
	errChan := make(chan error, 1)
	go func() {
		if err := ctnr.Router.Run(fmt.Sprintf(":%d", ctnr.Variable.Api.Port)); err != nil {
			errChan <- fmt.Errorf("Router failed to run :%w", err)
		}
	}()
	return <-errChan
}
