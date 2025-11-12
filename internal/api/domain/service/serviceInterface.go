package service

type IService interface {
	//voice
	GetAudioChunk(id string, idx int32) ([]byte, error)
	GetCount(id string) (int32, error)

	CreateText(id string, speechIndex int32, entityIdx int, totalDuration, duration float64, chunk []byte) (string, error)
	CalculateSlienceDuration(chunk []byte) float64
	GetWavDuration(audio []byte) (float64, error)
	RemoveAll(dir string) error

	CompleteConversion(Id string, count int)
	UpdateContent(id string, content string) error
}
