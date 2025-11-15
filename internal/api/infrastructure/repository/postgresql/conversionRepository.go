package postgresql

import (
	sharedModel "github.com/daeroworld/shared/model"
	"gorm.io/gorm/clause"

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

func (repo *ConversionRepository) Create(f *sharedModel.TextConversion) (*sharedModel.TextConversion, error) {
	if err := repo.postgresql.Driver.Create(f).Error; err != nil {
		return nil, err
	}
	return f, nil
}

func (repo *ConversionRepository) Upsert(tc *sharedModel.TextConversion) (*sharedModel.TextConversion, error) {
	if err := repo.postgresql.Driver.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(tc).Error; err != nil {
		return nil, err
	}
	return tc, nil
}
