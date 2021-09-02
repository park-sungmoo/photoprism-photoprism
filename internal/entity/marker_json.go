package entity

import (
	"encoding/json"
	"time"
)

// MarshalJSON returns the JSON encoding.
func (m *Marker) MarshalJSON() ([]byte, error) {
	var subj *Subject
	var name string

	if subj = m.Subject(); subj == nil {
		name = m.MarkerName
	} else {
		name = subj.SubjectName
	}

	return json.Marshal(&struct {
		UID        string
		FileUID    string
		Type       string
		Src        string
		Name       string
		X          float32
		Y          float32
		W          float32
		H          float32
		SubjectUID string `json:",omitempty"`
		SubjectSrc string `json:",omitempty"`
		FaceID     string `json:",omitempty"`
		FaceThumb  string `json:",omitempty"`
		Size       int    `json:",omitempty"`
		Score      int    `json:",omitempty"`
		Review     bool   `json:",omitempty"`
		Invalid    bool   `json:",omitempty"`
		CreatedAt  time.Time
	}{
		UID:        m.MarkerUID,
		FileUID:    m.FileUID,
		Type:       m.MarkerType,
		Src:        m.MarkerSrc,
		Name:       name,
		X:          m.X,
		Y:          m.Y,
		W:          m.W,
		H:          m.H,
		SubjectUID: m.SubjectUID,
		SubjectSrc: m.SubjectSrc,
		FaceID:     m.FaceID,
		FaceThumb:  m.FaceThumb,
		Size:       m.Size,
		Score:      m.Score,
		Review:     m.Review,
		Invalid:    m.MarkerInvalid,
		CreatedAt:  m.CreatedAt,
	})
}
