package controller

import (
	api "github.com/daeroworld/shared/api"
)

type IController interface {
	api.IController
	Create(id string) error
}
