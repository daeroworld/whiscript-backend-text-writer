package controller

import (
	"text/writer/internal/api/domain/service"
)

type Controller struct {
	svc service.IService
}

func NewController(svc service.IService) *Controller {
	return &Controller{
		svc: svc,
	}
}

func (ctrl *Controller) Create(id string) (int32, error) {
	cnt, err := ctrl.svc.GetCount(id)
	if err != nil {
		return 0, err
	}
	totalDuration := 0.0
	entityIdx := 0
	for idx := range cnt {
		chunk, err := ctrl.svc.GetAudioChunk(id, idx)
		if err != nil {
			continue
		}
		duration, err := ctrl.svc.CreateText(id, idx, &entityIdx, chunk, totalDuration)
		if err != nil {
			continue
		}
		totalDuration += duration
	}
	ctrl.svc.CompleteConversion(id, entityIdx)
	return int32(entityIdx), nil
}
