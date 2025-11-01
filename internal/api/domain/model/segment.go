package model

type WhisperSegment struct {
	ID    int     `json:"id"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text"`
}

type WhisperJson struct {
	Segments []WhisperSegment `json:"segments"`
}
