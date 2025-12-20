package main

import (
	"fmt"
	"log"
	"text/writer/configuration"
)

func main() {
	log.Printf("starting app")
	startContainer()
}
func startContainer() error {
	ctnr := configuration.NewContainer()
	errChan := make(chan error, 1)
	log.Printf("gRPC server listening on %s:%d\n", ctnr.Variable.Api.Ip, ctnr.Variable.Api.Port)
	go func() {
		if err := ctnr.DefineGrpc(); err != nil {
			errChan <- fmt.Errorf("gRPC server failed: %w", err)
		}
	}()

	return <-errChan
}
