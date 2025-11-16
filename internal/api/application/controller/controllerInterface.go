package controller

import (
	api "github.com/daeroworld/shared/api"
)

type IController interface {
	api.IController
	Create(id string) (int32, error)
	CreateSync(id string) (int32, error)
	UpdateContent(id, content string) error
}
