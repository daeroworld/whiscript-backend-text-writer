package service

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"os"
	"path/filepath"
	"text/writer/internal/api/domain/business"
	"text/writer/internal/api/domain/model"
	"text/writer/internal/api/infrastructure/client/voice"
	"text/writer/internal/api/infrastructure/repository/postgresql"

	sharedModel "github.com/daeroworld/shared/model"
)

type Service struct {
	biz            business.IBusiness
	voiceClnt      voice.IVoiceReader
	textRepo       postgresql.ITextRepository
	conversionRepo postgresql.IConversionRepository
}

func NewService(whisperBiz business.IBusiness, voiceClnt voice.IVoiceReader, textRepo postgresql.ITextRepository, conversionRepo postgresql.IConversionRepository) *Service {
	return &Service{
		biz:            whisperBiz,
		voiceClnt:      voiceClnt,
		textRepo:       textRepo,
		conversionRepo: conversionRepo,
	}
}

func (svc *Service) Init(filename string) (*sharedModel.TextConversion, error) {
	return svc.conversionRepo.Create(sharedModel.CreateTextConversion(filename, 0))
}

func (svc *Service) GetChunkCount(id string) (int32, error) {
	res, err := svc.voiceClnt.Count(id)
	if err != nil {
		return 0, err
	}
	return res.GetCount(), err
}

func (svc *Service) GetAudioChunk(id string, idx int32) ([]byte, error) {
	chunkRes, err := svc.voiceClnt.Retrieve(id, idx)
	if err != nil {
		return nil, err
	}
	return chunkRes.GetChunk(), nil
}

func (svc *Service) CreateText(id uint, chunkIndex int, totalDuration, duration float64, chunk []byte) (string, int, error) {
	filePath, err := svc.writeTempAudio(id, chunk)
	if err != nil {
		return "", 0, err
	}
	jsonPath, err := svc.run(filePath)
	segments, err := svc.load(jsonPath)
	createdCnt := 0
	for sentenceIndex, seg := range segments {

		for wordIdx, word := range seg.Words {
			StartFromTotal := totalDuration + word.Start
			EndFromTotal := totalDuration + word.End
			svc.textRepo.Create(sharedModel.CreateText(id, chunkIndex, sentenceIndex, wordIdx, StartFromTotal, EndFromTotal, word.Word))
			createdCnt++
		}

	}

	return filepath.Dir(filePath), createdCnt, nil
}

func (svc *Service) CalculateSlienceDuration(chunk []byte) float64 {
	bits := binary.LittleEndian.Uint64(chunk)
	return math.Float64frombits(bits)
}

func (svc *Service) GetWavDuration(audio []byte) (float64, error) {

	r := bytes.NewReader(audio)

	// ---- WAV Header ----
	// Byte 22-23: NumChannels (uint16)
	if len(audio) < 22 {
		return 0, nil
	}
	r.Seek(22, 0)
	var numChannels uint16
	if err := binary.Read(r, binary.LittleEndian, &numChannels); err != nil {
		return 0, err
	}

	// Byte 24-27: SampleRate (uint32)
	var sampleRate uint32
	if err := binary.Read(r, binary.LittleEndian, &sampleRate); err != nil {
		return 0, err
	}

	// Byte 34-35: BitsPerSample (uint16)
	r.Seek(34, 0)
	var bitsPerSample uint16
	if err := binary.Read(r, binary.LittleEndian, &bitsPerSample); err != nil {
		return 0, err
	}

	bytesPerSample := bitsPerSample / 8

	// ---- Find the "data" chunk dynamically ----
	pos := int64(12) // skip RIFF + WAVE header
	for {
		r.Seek(pos, 0)

		var chunkID [4]byte
		if err := binary.Read(r, binary.LittleEndian, &chunkID); err != nil {
			return 0, err
		}

		var chunkSize uint32
		if err := binary.Read(r, binary.LittleEndian, &chunkSize); err != nil {
			return 0, err
		}

		if string(chunkID[:]) == "data" {
			// We found the actual audio PCM data
			dataSize := chunkSize
			samples := float64(dataSize) / float64(bytesPerSample)
			durationSeconds := samples / float64(numChannels) / float64(sampleRate)
			return durationSeconds, nil
		}

		// Move to next chunk
		pos += int64(chunkSize) + 8
		if pos >= int64(len(audio)) {
			return 0, errors.New("data chunk not found in WAV file")
		}
	}
}

func (svc *Service) writeTempAudio(id uint, audio []byte) (string, error) {
	tempDir, err := os.MkdirTemp("", "audio")
	if err != nil {
		return "", err
	}

	// Generate filename using Id
	filePath := filepath.Join(tempDir, string(id)+".wav")

	// Create file
	f, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Write audio bytes
	if _, err := f.Write(audio); err != nil {
		return "", err
	}

	return filePath, nil
}

func (svc *Service) run(filename string) (string, error) {
	return svc.biz.Run(filename)
}

func (svc *Service) load(path string) ([]model.Segment, error) {
	return svc.biz.Load(path)
}

func (svc *Service) RemoveAll(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		// we already have the data, so return data + error
		return err
	}
	return nil
}

func (svc *Service) CompleteConversion(tc *sharedModel.TextConversion, count int) {
	tc.Count = count
	svc.conversionRepo.Upsert(tc)
}

func (svc *Service) UpdateContent(id string, content string) error {
	return svc.textRepo.UpdateContent(id, content)
}
