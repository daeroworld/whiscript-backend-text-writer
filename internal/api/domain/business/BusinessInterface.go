package business

import "text/writer/internal/api/domain/model"

type IBusiness interface {
	Run(filePath string) (json string, err error)
	Load(filePath string) ([]model.Segment, error)
}
