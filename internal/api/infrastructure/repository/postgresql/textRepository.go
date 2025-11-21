package postgresql

import (
	"fmt"

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

func (repo *TextRepository) SortSentenceIndex(fileId uint, indexSpace int) error {
	sql := fmt.Sprintf(`
    WITH ordered AS (
        SELECT
            id,
            ROW_NUMBER() OVER (PARTITION BY file_id ORDER BY chunk, sentence, word) AS rn
        FROM text
        WHERE file_id = ?
    )
    UPDATE text t
    SET sentence = o.rn * %d
    FROM ordered o
    WHERE t.id = o.id;
    `, indexSpace) // inject indexSpace here

	return repo.postgresql.Driver.Exec(sql, fileId).Error
}
