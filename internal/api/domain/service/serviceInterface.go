package service

type IService interface {
	//voice
	GetAudioChunk(id string, idx int32) ([]byte, error)
	GetCount(id string) (int32, error)

	//text
	CreateText(Id string, idx int32, entityIdx *int, audioChunk []byte, totalDuration float64) (duration float64, err error)
	CompleteConversion(Id string, count int)
	UpdateContent(id string, content string) error
}
