package business

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/writer/internal/api/domain/model"
)

type WhisperBusiness struct {
	semaphore chan struct{}
}

func NewWhisperBusiness(maxWorkers int) *WhisperBusiness {
	return &WhisperBusiness{
		semaphore: make(chan struct{}, maxWorkers),
	}
}

func (wb *WhisperBusiness) Run(filename string) (string, error) {

	wb.semaphore <- struct{}{}
	defer func() { <-wb.semaphore }()

	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	jsonOut := base + ".json"
	dir := filepath.Dir(filename)
	file := filepath.Base(filename)
	cmd := exec.Command("whisper", file, "--model", "medium", "--output_format", "json", "--language", "Korean", "--word_timestamps", "True")
	cmd.Dir = dir // <-- IMPORTANT

	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))

	return jsonOut, err
}

func (wb *WhisperBusiness) Load(filePath string) ([]model.Segment, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var w model.Json
	err = json.Unmarshal(b, &w)
	return w.Segments, err
}
