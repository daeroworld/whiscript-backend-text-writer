package postgresql_test

import (
	"os"
	"testing"
	"time"

	"text/writer/configuration"
	"text/writer/internal/api/infrastructure/repository/postgresql"

	sharedModel "github.com/daeroworld/shared/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/daeroworld/shared/database"
)

var (
	db   *database.PostgresqlWrapper
	repo postgresql.ITextRepository
)

func setupTestDB() *database.PostgresqlWrapper {
	config := configuration.NewVariable()
	return database.ConnectPostgresqlDatabase(config.Database)
}

func TestMain(m *testing.M) {
	db = setupTestDB()
	repo = postgresql.NewTextRepository(db)
	code := m.Run()
	os.Exit(code)
}

func measure(t *testing.T, name string, fn func()) {
	start := time.Now()
	fn()
	t.Logf("%s took %v", name, time.Since(start))
}

func createTestFile() *sharedModel.File {
	return &sharedModel.File{
		Id:       uuid.New(),
		Filename: "TEST",
		Status:   int(sharedModel.PROCESSING),

		// Add other required fields based on your File model
	}
}

func TestSortSentenceIndex(t *testing.T) {

	// insert sample texts out of order
	const fileId = 5
	texts := []*sharedModel.Text{
		{FileId: fileId, Chunk: 1, Sentence: 2000, Word: 1, Content: "B"},
		{FileId: fileId, Chunk: 0, Sentence: 1000, Word: 0, Content: "A"},
		{FileId: fileId, Chunk: 0, Sentence: 2000, Word: 1, Content: "C"},
	}
	for _, t := range texts {
		repo.Create(t)
	}

	err := repo.SortSentenceIndex(fileId, 1000)
	assert.NoError(t, err)

	var sorted []sharedModel.Text
	db.Driver.Where("file_id = ?", fileId).Order("chunk, sentence, word").Find(&sorted)

	// Check that sentence numbers are renumbered as 1000, 2000, 3000
	assert.Equal(t, 1000, sorted[0].Sentence)
	assert.Equal(t, 2000, sorted[1].Sentence)
	assert.Equal(t, 3000, sorted[2].Sentence)
}
