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

type StableTsBusiness struct {
}

func NewStableTsBusiness() *StableTsBusiness {
	return &StableTsBusiness{}
}

func (wb *StableTsBusiness) Run(filename string) (string, error) {
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	jsonOut := base + ".json"
	dir := filepath.Dir(filename)
	file := filepath.Base(filename)
	cmd := exec.Command("stable-ts", file, "--model", "medium", "--output_format", "json", "--language", "Korean", "--word_timestamps")
	cmd.Dir = dir // <-- IMPORTANT

	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))

	return jsonOut, err
}

func (wb *StableTsBusiness) Load(filePath string) ([]model.Segment, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var w model.Json
	err = json.Unmarshal(b, &w)
	return w.Segments, err
}
