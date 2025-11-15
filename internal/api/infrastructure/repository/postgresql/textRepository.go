package postgresql

import (
	sharedModel "github.com/daeroworld/shared/model"

	"github.com/daeroworld/shared/database"
)

type TextRepository struct {
	postgresql *database.PostgresqlWrapper
}

func NewTextRepository(postgresql *database.PostgresqlWrapper) *TextRepository {
	return &TextRepository{
		postgresql: postgresql,
	}
}

func (repo *TextRepository) Create(f *sharedModel.Text) (*sharedModel.Text, error) {
	if err := repo.postgresql.Driver.Create(f).Error; err != nil {
		return nil, err
	}
	return f, nil
}

func (repo *TextRepository) UpdateContent(id string, content string) error {
	return repo.postgresql.Driver.
		Model(&sharedModel.Text{}).
		Where("id = ?", id).
		Update("edit_content", content).Error
}
