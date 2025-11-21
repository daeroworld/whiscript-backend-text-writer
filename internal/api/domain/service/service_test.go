package service_test

import (
	"os"
	"testing"
	"text/writer/internal/api/domain/model"
	"text/writer/internal/api/domain/service"

	sharedModel "github.com/daeroworld/shared/model"
	"github.com/daeroworld/shared/proto/voice"
	pb "github.com/daeroworld/shared/proto/voice"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockIndexBiz struct {
	mock.Mock
}

// ConvertSentenceIdex implements index.IIndexBusiness.
func (m *MockIndexBiz) ConvertSentenceIdex(int) int {
	panic("unimplemented")
}

// GetIndexSpace implements index.IIndexBusiness.
func (m *MockIndexBiz) GetIndexSpace() int {
	args := m.Called()
	return args.Int(0)
}

type MockWhisperBiz struct {
	mock.Mock
}

// Load implements business.IBusiness.
func (m *MockWhisperBiz) Load(filePath string) ([]model.Segment, error) {
	panic("unimplemented")
}

// Run implements business.IBusiness.
func (m *MockWhisperBiz) Run(filePath string) (json string, err error) {
	panic("unimplemented")
}

type MockVoiceClient struct {
	mock.Mock
}

func (m *MockVoiceClient) Count(id string) (*pb.VoiceCountResponse, error) {
	args := m.Called(id)
	return args.Get(0).(*pb.VoiceCountResponse), args.Error(1)
}

func (m *MockVoiceClient) Retrieve(id string, idx int32) (*pb.VoiceRetrieveResponse, error) {
	args := m.Called(id, idx)
	return args.Get(0).(*pb.VoiceRetrieveResponse), args.Error(1)
}

type MockTextRepo struct {
	mock.Mock
}

// CreateSentence implements postgresql.ITextRepository.
func (m *MockTextRepo) CreateSentence(t *sharedModel.Text, filename string) (*sharedModel.Text, error) {
	args := m.Called(t, filename)
	return args.Get(0).(*sharedModel.Text), args.Error(1)
}

// Update implements postgresql.ITextRepository.
func (m *MockTextRepo) Update(t *sharedModel.Text) (*sharedModel.Text, error) {
	args := m.Called(t)
	return args.Get(0).(*sharedModel.Text), args.Error(1)
}

// Create implements postgresql.ITextRepository.
func (m *MockTextRepo) Create(f *sharedModel.Text) (*sharedModel.Text, error) {
	panic("unimplemented")
}

// SortSentenceIndex implements postgresql.ITextRepository.
func (m *MockTextRepo) SortSentenceIndex(fileId uint, indexSpace int) error {
	panic("unimplemented")
}

func (m *MockTextRepo) Save(text *sharedModel.Text) (*sharedModel.Text, error) {
	args := m.Called(text)
	return text, args.Error(0)
}

func (m *MockTextRepo) UpdateContent(id string, content string) error {
	args := m.Called(id, content)
	return args.Error(0)
}

type MockConversionRepo struct {
	mock.Mock
}

// Create implements postgresql.IConversionRepository.
func (m *MockConversionRepo) Create(f *sharedModel.TextConversion) (*sharedModel.TextConversion, error) {
	panic("unimplemented")
}

// Upsert implements postgresql.IConversionRepository.
func (m *MockConversionRepo) Upsert(e *sharedModel.TextConversion) (*sharedModel.TextConversion, error) {
	panic("unimplemented")
}

func (m *MockConversionRepo) Save(conv *sharedModel.TextConversion) (*sharedModel.TextConversion, error) {
	args := m.Called(conv)
	return conv, args.Error(0)
}

var (
	mockIdxBiz     *MockIndexBiz
	mockWhisperBiz *MockWhisperBiz
	mockVoiceClnt  *MockVoiceClient
	mockTextRepo   *MockTextRepo
	mockConvRepo   *MockConversionRepo
	svc            *service.Service
)

func TestMain(m *testing.M) {
	mockIdxBiz = new(MockIndexBiz)
	mockWhisperBiz = new(MockWhisperBiz)
	mockVoiceClnt = new(MockVoiceClient)
	mockTextRepo = new(MockTextRepo)
	mockConvRepo = new(MockConversionRepo)

	svc = service.NewService(mockIdxBiz, mockWhisperBiz, mockVoiceClnt, mockTextRepo, mockConvRepo)
	code := m.Run()
	os.Exit(code)
}

type MockChunkResponse struct {
	chunk []byte
}

func (r *MockChunkResponse) GetChunk() []byte { return r.chunk }

type MockCountResponse struct {
	count int32
}

func (r *MockCountResponse) GetCount() int32 { return r.count }

// ---- 실제 테스트 ----

func TestGetChunkCount(t *testing.T) {
	mockVoiceClnt = new(MockVoiceClient)
	svc = service.NewService(mockIdxBiz, mockWhisperBiz, mockVoiceClnt, mockTextRepo, mockConvRepo)

	mockVoiceClnt.On("Count", "abc123").Return(&voice.VoiceCountResponse{Count: 5}, nil)

	count, err := svc.GetChunkCount("abc123")

	assert.NoError(t, err)
	assert.Equal(t, int32(5), count)
	mockVoiceClnt.AssertExpectations(t)
}

func TestPut(t *testing.T) {
	mockRepo := new(MockTextRepo)
	mockIdx := new(MockIndexBiz)

	svc = service.NewService(mockIdx, mockWhisperBiz, mockVoiceClnt, mockRepo, mockConvRepo)

	id := uuid.New().String()
	const filename = "zzvideo1.mp4"
	const sentence = 1656
	word := 1
	start := 0.5
	end := 2.0
	content := "ccHello world"

	text := sharedModel.CreateTextForUpsert(id, sentence, word, start, end, content)

	// Expect GetIndexSpace to be called
	mockIdx.On("GetIndexSpace").Return(1000)

	// Expect Update to be called because sentence % indexSpace == 0
	mockRepo.On("Update", text).Return(text, nil)

	// Expect CreateSentence to be called
	mockRepo.On("CreateSentence", text, filename).Return(text, nil)

	// Call the method
	generatedText, err := svc.Put(id, filename, sentence, word, start, end, content)

	assert.NoError(t, err)
	assert.Equal(t, text, generatedText)
}
