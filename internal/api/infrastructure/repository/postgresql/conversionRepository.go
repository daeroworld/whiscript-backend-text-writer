package postgresql

import (
	sharedModel "github.com/daeroworld/shared/model"

	"github.com/daeroworld/shared/database"
)

type ConversionRepository struct {
	postgresql *database.PostgresqlWrapper
}

func NewConversionRepository(postgresql *database.PostgresqlWrapper) *ConversionRepository {
	return &ConversionRepository{
		postgresql: postgresql,
	}
}

func (repo *ConversionRepository) Save(f *sharedModel.TextConversion) (*sharedModel.TextConversion, error) {
	if err := repo.postgresql.Driver.Create(f).Error; err != nil {
		return nil, err
	}
	return f, nil
}

