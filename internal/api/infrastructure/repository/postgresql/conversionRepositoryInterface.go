package postgresql

import (
	"text/writer/internal/api/infrastructure/repository"

	sharedModel "github.com/daeroworld/shared/model"
)

type IConversionRepository interface {
	repository.IRepository[sharedModel.TextConversion]
}
