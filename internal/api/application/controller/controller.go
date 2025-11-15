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
	tcv, err := ctrl.svc.Init(id)
	if err != nil {
		return 0, err
	}

	cnt, err := ctrl.svc.GetChunkCount(id)
	if err != nil {
		return 0, err
	}

	handleSilence := func(chunk []byte) (float64, error) {
		return ctrl.svc.CalculateSlienceDuration(chunk), nil
	}

	handleSpeech := func(id string, chunkIdx int32, entityIdx int, totalDuration, duration float64, chunk []byte) error {
		dir, err := ctrl.svc.CreateText(id, chunkIdx, entityIdx, totalDuration, duration, chunk)
		if err != nil {
			return err
		}
		go func(dir string) {
			ctrl.svc.RemoveAll(dir)
		}(dir)
		return nil
	}

	totalDuration := 0.0
	entityCount := 0
	for chunkIdx := range cnt {
		chunk, err := ctrl.svc.GetAudioChunk(id, chunkIdx)
		if err != nil {
			continue
		}
		duration, err := ctrl.svc.GetWavDuration(chunk)
		if err != nil {
			continue
		}
		if duration == 0 { // if it is silence
			duration, err = handleSilence(chunk)
			totalDuration += duration
			continue
		}
		err = handleSpeech(id, chunkIdx, entityCount, totalDuration, duration, chunk)
		if err != nil {
			continue
		}
		totalDuration += duration
		entityCount += cnt
	}
	ctrl.svc.CompleteConversion(tcv, entityCount)
	return int32(entityCount), nil
}

func (ctrl *Controller) UpdateContent(id, content string) error {
	return ctrl.svc.UpdateContent(id, content)
}
