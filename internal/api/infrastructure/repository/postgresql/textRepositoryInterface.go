package postgresql

import (
	sharedModel "github.com/daeroworld/shared/model"
)

type ITextRepository interface {
	Create(f *sharedModel.Text) (*sharedModel.Text, error)
	Update(t *sharedModel.Text) (*sharedModel.Text, error)
	SortSentenceIndex(fileId uint, indexSpace int) error
	CreateSentence(t *sharedModel.Text, filename string) (*sharedModel.Text, error)
	UpdateContent(id string, content string) error
}
