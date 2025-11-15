package postgresql

import (
	sharedModel "github.com/daeroworld/shared/model"
)

type ITextRepository interface {
	Create(f *sharedModel.Text) (*sharedModel.Text, error)
	UpdateContent(id string, content string) error
}
