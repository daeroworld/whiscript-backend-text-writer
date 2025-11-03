package service

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/writer/internal/api/domain/model"
	"text/writer/internal/api/infrastructure/client/voice"
	"text/writer/internal/api/infrastructure/repository"
	"text/writer/internal/api/infrastructure/repository/postgresql"

	sharedModel "github.com/daeroworld/shared/model"
)

type Service struct {
	voiceClnt      *voice.VoiceReaderClient
	textRepo       *postgresql.TextRepository
	conversionRepo repository.IRepository[sharedModel.TextConversion]
}

func NewService(voiceClnt *voice.VoiceReaderClient, textRepo *postgresql.TextRepository, conversionRepo repository.IRepository[sharedModel.TextConversion]) *Service {
	return &Service{
		voiceClnt:      voiceClnt,
		textRepo:       textRepo,
		conversionRepo: conversionRepo,
	}
}

func (svc *Service) GetCount(id string) (int32, error) {
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

func (svc *Service) CreateText(id string, voiceIdx int32, entityIdx *int, audioChunk []byte, totalDuration float64) (float64, error) {
	duration, err := svc.getWavDuration(audioChunk)
	if err != nil {
		return 0, err
	}
	if duration > 0 {
		svc.createSpeechText(id, voiceIdx, entityIdx, totalDuration, duration, audioChunk)
		return duration, nil
	}
	svc.calculateSlienceDuration(&duration, audioChunk)

	return duration, nil
}

func (svc *Service) createSpeechText(id string, speechIndex int32, entityIdx *int, totalDuration, duration float64, chunk []byte) error {
	filePath, err := svc.writeTempAudio(id, chunk)
	if err != nil {
		return err
	}
	jsonPath, err := svc.runWhisper(filePath, chunk)
	segments, err := svc.loadWhisperJson(jsonPath)

	for textIndex, seg := range segments {
		StartFromTotal := totalDuration + seg.Start
		EndFromTotal := totalDuration + seg.End
		*entityIdx = *entityIdx + 1
		svc.textRepo.Save(sharedModel.CreateText(id, int(speechIndex), textIndex, *entityIdx, StartFromTotal, EndFromTotal, seg.Text))
	}
	return nil
}

func (svc *Service) createSilenceText(id string, speechIndex int32, entityIdx *int, totalDuration float64, duration *float64, chunk []byte) error {
	BytesToFloat64 := func(b []byte) float64 {
		bits := binary.LittleEndian.Uint64(b)
		return math.Float64frombits(bits)
	}
	*duration = BytesToFloat64(chunk)
	StartFromTotal := totalDuration
	EndFromTotal := totalDuration + *duration
	*entityIdx = *entityIdx + 1
	svc.textRepo.Save(sharedModel.CreateText(id, int(speechIndex), *entityIdx, 0, StartFromTotal, EndFromTotal, "."))
	return nil
}

func (svc *Service) calculateSlienceDuration(duration *float64, chunk []byte) error {
	bits := binary.LittleEndian.Uint64(chunk)
	*duration = math.Float64frombits(bits)
	return nil
}

func (svc *Service) getWavDuration(audio []byte) (float64, error) {

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

func (svc *Service) writeTempAudio(id string, audio []byte) (string, error) {
	tempDir, err := os.MkdirTemp("", "audio")
	if err != nil {
		return "", err
	}

	// Generate filename using ID
	filePath := filepath.Join(tempDir, id+".wav")

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

func (svc *Service) runWhisper(filename string, chunk []byte) (string, error) {
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	jsonOut := base + ".json"
	dir := filepath.Dir(filename)
	file := filepath.Base(filename)
	cmd := exec.Command("whisper", file, "--model", "medium", "--output_format", "json", "--language", "Korean")
	cmd.Dir = dir // <-- IMPORTANT

	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))

	return jsonOut, err
}

func (svc *Service) loadWhisperJson(path string) ([]model.WhisperSegment, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var w model.WhisperJson
	err = json.Unmarshal(b, &w)
	return w.Segments, err
}

func (svc *Service) sliceAudio(filePath string, seg model.WhisperSegment) ([]byte, error) {
	dir := filepath.Dir(filePath)
	out := filepath.Join(dir, fmt.Sprintf("chunk_%03d.wav", seg.ID))
	cmd := exec.Command("ffmpeg",
		"-y",
		"-i", filePath,
		"-ss", fmt.Sprintf("%.3f", seg.Start),
		"-to", fmt.Sprintf("%.3f", seg.End),
		"-c", "copy",
		out,
	)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// Return file data as []byte
	return os.ReadFile(out)
}

func (svc *Service) CompleteConversion(id string, count int) {
	svc.conversionRepo.Save(sharedModel.CreateTextConversion(id, count))
}


