package entity

type SubjectMap map[string]Subject

func (m SubjectMap) Get(name string) Subject {
	if result, ok := m[name]; ok {
		return result
	}

	return Subject{}
}

func (m SubjectMap) Pointer(name string) *Subject {
	if result, ok := m[name]; ok {
		return &result
	}

	return &Subject{}
}

var SubjectFixtures = SubjectMap{
	"john-doe": Subject{
		SubjectUID:   "jqu0xs11qekk9jx8",
		SubjectSlug:  "john-doe",
		SubjectName:  "John Doe",
		SubjectType:  SubjectPerson,
		SubjectSrc:   SrcManual,
		Favorite:     true,
		Private:      false,
		Excluded:     false,
		SubjectBio:   "Subject Description",
		SubjectNotes: "Short Note",
		MetadataJSON: []byte(""),
		FileCount:    1,
		CreatedAt:    TimeStamp(),
		UpdatedAt:    TimeStamp(),
		DeletedAt:    nil,
	},
	"joe-biden": Subject{
		SubjectUID:   "jqy3y652h8njw0sx",
		SubjectSlug:  "joe-biden",
		SubjectName:  "Joe Biden",
		SubjectType:  SubjectPerson,
		SubjectSrc:   SrcMarker,
		Favorite:     false,
		Private:      false,
		Excluded:     false,
		SubjectBio:   "",
		SubjectNotes: "",
		MetadataJSON: []byte(""),
		FileCount:    1,
		CreatedAt:    TimeStamp(),
		UpdatedAt:    TimeStamp(),
		DeletedAt:    nil,
	},
	"dangling": Subject{
		SubjectUID:   "jqy1y111h1njaaaa",
		SubjectSlug:  "dangling-subject",
		SubjectName:  "Dangling Subject",
		SubjectType:  SubjectPerson,
		SubjectSrc:   SrcMarker,
		Favorite:     false,
		Private:      false,
		Excluded:     false,
		SubjectBio:   "",
		SubjectNotes: "",
		MetadataJSON: []byte(""),
		FileCount:    0,
		CreatedAt:    TimeStamp(),
		UpdatedAt:    TimeStamp(),
		DeletedAt:    nil,
	},
	"jane-doe": Subject{
		SubjectUID:   "jqy1y111h1njaaab",
		SubjectSlug:  "jane-doe",
		SubjectName:  "Jane Doe",
		SubjectType:  SubjectPerson,
		SubjectSrc:   SrcMarker,
		Favorite:     false,
		Private:      false,
		Excluded:     false,
		SubjectBio:   "",
		SubjectNotes: "",
		MetadataJSON: []byte(""),
		FileCount:    3,
		CreatedAt:    TimeStamp(),
		UpdatedAt:    TimeStamp(),
		DeletedAt:    nil,
	},
	"actress-1": Subject{
		SubjectUID:   "jqy1y111h1njaaac",
		SubjectSlug:  "actress-a",
		SubjectName:  "Actress A",
		SubjectType:  SubjectPerson,
		SubjectSrc:   SrcMarker,
		Favorite:     false,
		Private:      false,
		SubjectNotes: "",
		MetadataJSON: []byte(""),
		CreatedAt:    TimeStamp(),
		UpdatedAt:    TimeStamp(),
		DeletedAt:    nil,
	},
	"actor-1": Subject{
		SubjectUID:   "jqy1y111h1njaaad",
		SubjectSlug:  "actor-a",
		SubjectName:  "Actor A",
		SubjectType:  SubjectPerson,
		SubjectSrc:   SrcMarker,
		Favorite:     false,
		Private:      false,
		SubjectNotes: "",
		MetadataJSON: []byte(""),
		CreatedAt:    TimeStamp(),
		UpdatedAt:    TimeStamp(),
		DeletedAt:    nil,
	},
}

// CreateSubjectFixtures inserts known entities into the database for testing.
func CreateSubjectFixtures() {
	for _, entity := range SubjectFixtures {
		Db().Create(&entity)
	}
}
