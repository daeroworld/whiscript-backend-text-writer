package postgresql

import (
	"fmt"

	sharedModel "github.com/daeroworld/shared/model"
	"gorm.io/gorm/clause"

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

func (repo *TextRepository) Update(t *sharedModel.Text) (*sharedModel.Text, error) {
	err := repo.postgresql.Driver.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"sentence":     t.Sentence,
				"word":         t.Word,
				"start":        t.Start,
				"end":          t.End,
				"edit_content": t.EditContent,
			}),
		}).
		Create(t).Error

	if err != nil {
		return nil, err
	}
	return t, nil
}

func (repo *TextRepository) CreateSentence(t *sharedModel.Text, filename string) (*sharedModel.Text, error) {
	sql := `
		INSERT INTO text (id, file_id, chunk, sentence, word, start, "end", content, edit_content, created_at)
		SELECT ?, tc.id, ?, ?, ?, ?, ?, ?, ?, ?
		FROM text_conversion tc
		WHERE tc.filename = ?;
	`

	err := repo.postgresql.Driver.Exec(
		sql,
		t.Id,
		t.Chunk,
		t.Sentence,
		t.Word,
		t.Start,
		t.End,
		t.Content,
		t.EditContent,
		t.CreatedAt,
		filename,
	).Error

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (repo *TextRepository) UpdateContent(id string, content string) error {
	return repo.postgresql.Driver.
		Model(&sharedModel.Text{}).
		Where("id = ?", id).
		Update("edit_content", content).Error
}

func (repo *TextRepository) SortSentenceIndex(fileId uint, indexSpace int) error {
	sql := fmt.Sprintf(`
WITH unique_sentences AS (
    SELECT
        file_id,
        chunk,
        sentence,
        ROW_NUMBER() OVER (PARTITION BY file_id ORDER BY chunk, sentence) AS new_sentence_idx
    FROM (
        SELECT DISTINCT file_id, chunk, sentence
        FROM text
        WHERE file_id = ?
    ) s
)
UPDATE text t
SET sentence = u.new_sentence_idx * %d
FROM unique_sentences u
WHERE t.file_id = u.file_id
  AND t.chunk = u.chunk
  AND t.sentence = u.sentence;
    `, indexSpace)

	return repo.postgresql.Driver.Exec(sql, fileId).Error
}
