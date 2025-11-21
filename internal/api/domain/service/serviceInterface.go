package service

import sharedModel "github.com/daeroworld/shared/model"

type IService interface {
	Init(filename string) (*sharedModel.TextConversion, error)
	//voice
	GetAudioChunk(id string, idx int32) ([]byte, error)
	GetChunkCount(id string) (int32, error)

	CreateText(id uint, speechIndex int, totalDuration, duration float64, chunk []byte) (string, int, error)
	CalculateSlienceDuration(chunk []byte) float64
	GetWavDuration(audio []byte) (float64, error)
	RemoveAll(dir string) error

	SortSentenceIndex(id uint) error
	CompleteConversion(tc *sharedModel.TextConversion, count int)
	Put(id, filename string, sentence, word int, start, end float64, content string) (*sharedModel.Text, error)
}
