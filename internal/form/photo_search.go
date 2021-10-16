package form

import (
	"time"
)

// PhotoSearch represents search form fields for "/api/v1/photos".
type PhotoSearch struct {
	Query     string    `form:"q"`
	Filter    string    `form:"filter"`
	ID        string    `form:"id"`
	Type      string    `form:"type"`
	Path      string    `form:"path"`
	Folder    string    `form:"folder"` // Alias for Path
	Name      string    `form:"name"`
	Filename  string    `form:"filename"`
	Original  string    `form:"original"`
	Title     string    `form:"title"`
	Hash      string    `form:"hash"`
	Primary   bool      `form:"primary"`
	Stack     bool      `form:"stack"`
	Unstacked bool      `form:"unstacked"`
	Stackable bool      `form:"stackable"`
	Video     bool      `form:"video"`
	Photo     bool      `form:"photo"`
	Raw       bool      `form:"raw"`
	Live      bool      `form:"live"`
	Scan      bool      `form:"scan"`
	Panorama  bool      `form:"panorama"`
	Error     bool      `form:"error"`
	Hidden    bool      `form:"hidden"`
	Archived  bool      `form:"archived"`
	Public    bool      `form:"public"`
	Private   bool      `form:"private"`
	Favorite  bool      `form:"favorite"`
	Unsorted  bool      `form:"unsorted"`
	Lat       float32   `form:"lat"`
	Lng       float32   `form:"lng"`
	Dist      uint      `form:"dist"`
	Fmin      float32   `form:"fmin"`
	Fmax      float32   `form:"fmax"`
	Chroma    uint8     `form:"chroma"`
	Diff      uint32    `form:"diff"`
	Mono      bool      `form:"mono"`
	Portrait  bool      `form:"portrait"`
	Geo       bool      `form:"geo"`
	Keywords  string    `form:"keywords"`
	Label     string    `form:"label"`
	Category  string    `form:"category"` // Moments
	Country   string    `form:"country"`  // Moments
	State     string    `form:"state"`    // Moments
	Year      string    `form:"year"`     // Moments
	Month     string    `form:"month"`    // Moments
	Day       string    `form:"day"`      // Moments
	Face      string    `form:"face"`     // UIDs
	Subject   string    `form:"subject"`  // UIDs
	Person    string    `form:"person"`   // Alias for Subject
	Subjects  string    `form:"subjects"` // Text
	People    string    `form:"people"`   // Alias for Subjects
	Album     string    `form:"album"`    // UIDs
	Albums    string    `form:"albums"`   // Text
	Color     string    `form:"color"`
	Faces     string    `form:"faces"` // Find or exclude faces if detected.
	Quality   int       `form:"quality"`
	Review    bool      `form:"review"`
	Camera    int       `form:"camera"`
	Lens      int       `form:"lens"`
	Before    time.Time `form:"before" time_format:"2006-01-02"`
	After     time.Time `form:"after" time_format:"2006-01-02"`
	Count     int       `form:"count" binding:"required" serialize:"-"`
	Offset    int       `form:"offset" serialize:"-"`
	Order     string    `form:"order" serialize:"-"`
	Merged    bool      `form:"merged" serialize:"-"`
}

func (f *PhotoSearch) GetQuery() string {
	return f.Query
}

func (f *PhotoSearch) SetQuery(q string) {
	f.Query = q
}

func (f *PhotoSearch) ParseQueryString() error {
	if err := ParseQueryString(f); err != nil {
		return err
	}

	if f.Path == "" && f.Folder != "" {
		f.Path = f.Folder
		f.Folder = ""
	}

	if f.Subject == "" && f.Person != "" {
		f.Subject = f.Person
		f.Person = ""
	}

	if f.Subjects == "" && f.People != "" {
		f.Subjects = f.People
		f.People = ""
	}

	if f.Filter != "" {
		if err := Unserialize(f, f.Filter); err != nil {
			return err
		}
	}

	return nil
}

// Serialize returns a string containing non-empty fields and values of a struct.
func (f *PhotoSearch) Serialize() string {
	return Serialize(f, false)
}

// SerializeAll returns a string containing all non-empty fields and values of a struct.
func (f *PhotoSearch) SerializeAll() string {
	return Serialize(f, true)
}

func NewPhotoSearch(query string) PhotoSearch {
	return PhotoSearch{Query: query}
}
