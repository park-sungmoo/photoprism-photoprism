package entity

import (
	"crypto/sha1"
	"encoding/base32"
	"encoding/json"
	"sync"
	"time"
)

var faceMutex = sync.Mutex{}

// Faces represents a Face slice.
type Faces []Face

// Face represents the face of a Subject.
type Face struct {
	ID            string          `gorm:"type:VARBINARY(42);primary_key;auto_increment:false;" json:"ID" yaml:"ID"`
	FaceSrc       string          `gorm:"type:VARBINARY(8);" json:"Src" yaml:"Src,omitempty"`
	SubjectUID    string          `gorm:"type:VARBINARY(42);index;" json:"SubjectUID" yaml:"SubjectUID,omitempty"`
	Collisions    int             `json:"Collisions" yaml:"Collisions,omitempty"`
	Samples       int             `json:"Samples" yaml:"Samples,omitempty"`
	Radius        float64         `json:"Radius" yaml:"Radius,omitempty"`
	EmbeddingJSON json.RawMessage `gorm:"type:MEDIUMBLOB;" json:"-" yaml:"EmbeddingJSON,omitempty"`
	CreatedAt     time.Time       `json:"CreatedAt" yaml:"CreatedAt,omitempty"`
	UpdatedAt     time.Time       `json:"UpdatedAt" yaml:"UpdatedAt,omitempty"`
	embedding     Embedding       `gorm:"-"`
}

// UnknownFace can be used as a placeholder for unknown faces.
var UnknownFace = Face{
	ID:            UnknownID,
	FaceSrc:       SrcDefault,
	SubjectUID:    UnknownPerson.SubjectUID,
	EmbeddingJSON: []byte{},
}

// CreateUnknownFace initializes the database with a placeholder for unknown faces.
func CreateUnknownFace() {
	_ = UnknownFace.Create()
}

// TableName returns the entity database table name.
func (Face) TableName() string {
	return "faces_dev4"
}

// NewFace returns a new face.
func NewFace(subjectUID, faceSrc string, embeddings Embeddings) *Face {
	result := &Face{
		SubjectUID: subjectUID,
		FaceSrc:    faceSrc,
	}

	if err := result.SetEmbeddings(embeddings); err != nil {
		log.Errorf("face: failed setting embeddings (%s)", err)
	}

	return result
}

// SetEmbeddings assigns face embeddings.
func (m *Face) SetEmbeddings(embeddings Embeddings) (err error) {
	m.embedding, m.Radius, m.Samples = EmbeddingsMidpoint(embeddings)
	m.EmbeddingJSON, err = json.Marshal(m.embedding)

	if err != nil {
		return err
	}

	s := sha1.Sum(m.EmbeddingJSON)
	m.ID = base32.StdEncoding.EncodeToString(s[:])
	m.UpdatedAt = Timestamp()

	if m.CreatedAt.IsZero() {
		m.CreatedAt = m.UpdatedAt
	}

	return nil
}

// Embedding returns parsed face embedding.
func (m *Face) Embedding() Embedding {
	if len(m.EmbeddingJSON) == 0 {
		return Embedding{}
	} else if len(m.embedding) > 0 {
		return m.embedding
	} else if err := json.Unmarshal(m.EmbeddingJSON, &m.embedding); err != nil {
		log.Errorf("failed parsing face embedding json: %s", err)
	}

	return m.embedding
}

// Save updates the existing or inserts a new face.
func (m *Face) Save() error {
	faceMutex.Lock()
	defer faceMutex.Unlock()

	return Save(m, "ID")
}

// Create inserts the face to the database.
func (m *Face) Create() error {
	faceMutex.Lock()
	defer faceMutex.Unlock()

	return Db().Create(m).Error
}

// Delete removes the face from the database.
func (m *Face) Delete() error {
	return Db().Delete(m).Error
}

// Update a face property in the database.
func (m *Face) Update(attr string, value interface{}) error {
	return UnscopedDb().Model(m).Update(attr, value).Error
}

// Updates face properties in the database.
func (m *Face) Updates(values interface{}) error {
	return UnscopedDb().Model(m).Updates(values).Error
}

// FirstOrCreateFace returns the existing entity, inserts a new entity or nil in case of errors.
func FirstOrCreateFace(m *Face) *Face {
	result := Face{}

	if err := UnscopedDb().Where("id = ?", m.ID).First(&result).Error; err == nil {
		return &result
	} else if createErr := m.Create(); createErr == nil {
		return m
	} else if err := UnscopedDb().Where("id = ?", m.ID).First(&result).Error; err == nil {
		return &result
	} else {
		log.Errorf("face: %s (find or create %s)", createErr, m.ID)
	}

	return nil
}

// FindFace returns an existing entity if exists.
func FindFace(id string) *Face {
	if id == "" {
		return nil
	}

	result := Face{}

	db := Db()
	db = db.Where("id = ?", id)

	if err := db.First(&result).Error; err != nil {
		return nil
	}

	return &result
}
