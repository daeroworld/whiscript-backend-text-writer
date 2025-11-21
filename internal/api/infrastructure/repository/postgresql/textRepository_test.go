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
	const fileId = 77
	texts := []*sharedModel.Text{
		{FileId: fileId, Chunk: 1, Sentence: 1000, Word: 1, Content: "Hello"},
		{FileId: fileId, Chunk: 1, Sentence: 1000, Word: 2, Content: "World"},
		// Chunk 1, sentence 2000
		{FileId: fileId, Chunk: 1, Sentence: 2000, Word: 1, Content: "Foo"},
		// Chunk 2, sentence 3000, 3 words (duplicate sentence)
		{FileId: fileId, Chunk: 2, Sentence: 3000, Word: 1, Content: "Bar"},
		{FileId: fileId, Chunk: 2, Sentence: 3000, Word: 2, Content: "Baz"},
		{FileId: fileId, Chunk: 2, Sentence: 3000, Word: 3, Content: "Qux"},
		// Chunk 2, sentence 4000
		{FileId: fileId, Chunk: 2, Sentence: 4000, Word: 1, Content: "End"},
	}
	for _, text := range texts {
		_, err := repo.Create(text)
		assert.NoError(t, err)
	}

	err := repo.SortSentenceIndex(fileId, 1000)
	assert.NoError(t, err)

	var sorted []sharedModel.Text
	err = db.Driver.Order("chunk, sentence, word").Where("file_id = ?", fileId).Find(&sorted).Error
	assert.NoError(t, err)

	expectedSentences := []int{
		1000, 1000, 2000, 3000, 3000,
		3000, 4000,
	}

	for i, s := range sorted {
		assert.Equal(t, expectedSentences[i], s.Sentence, "text at index %d", i)
	}
}

func TestCreateSentence(t *testing.T) {
	const filename = "vcvc.mp4"
	// 1. Insert a TextConversion
	tc := sharedModel.TextConversion{
		Filename: filename,
		Count:    10,
	}
	err := db.Driver.Create(&tc).Error
	assert.NoError(t, err)

	// 2. Create Text object (input to CreateSentence)
	input := &sharedModel.Text{
		Id:          uuid.New(),
		Chunk:       1,
		Sentence:    1000,
		Word:        1,
		Start:       0.5,
		End:         2.0,
		EditContent: "Hello world",
		CreatedAt:   time.Now(),
	}

	// 3. Call CreateSentence
	out, err := repo.CreateSentence(input, filename)
	assert.NoError(t, err)
	assert.NotNil(t, out)

	// 4. Verify inserted row in DB
	var saved sharedModel.Text
	err = db.Driver.Where("id = ?", input.Id).First(&saved).Error
	assert.NoError(t, err)

	// 5. Validate fields
	assert.Equal(t, input.Id, saved.Id)
	assert.Equal(t, tc.Id, saved.FileId)
	assert.Equal(t, input.Chunk, saved.Chunk)
	assert.Equal(t, input.Sentence, saved.Sentence)
	assert.Equal(t, input.Word, saved.Word)
	assert.Equal(t, input.Start, saved.Start)
	assert.Equal(t, input.End, saved.End)
	assert.Equal(t, input.Content, saved.Content)
	assert.Equal(t, input.EditContent, saved.EditContent)
}
