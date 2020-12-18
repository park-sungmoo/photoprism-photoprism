package entity

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/photoprism/photoprism/pkg/fs"
	"gopkg.in/yaml.v2"
)

var albumYamlMutex = sync.Mutex{}

// Yaml returns album data as YAML string.
func (m *Album) Yaml() ([]byte, error) {
	if err := Db().Model(m).Association("Photos").Find(&m.Photos).Error; err != nil {
		log.Errorf("album: %s (yaml)", err)
	}

	out, err := yaml.Marshal(m)

	if err != nil {
		return []byte{}, err
	}

	return out, err
}

// SaveAsYaml saves album data as YAML file.
func (m *Album) SaveAsYaml(fileName string) error {
	data, err := m.Yaml()

	if err != nil {
		return err
	}

	// Make sure directory exists.
	if err := os.MkdirAll(filepath.Dir(fileName), os.ModePerm); err != nil {
		return err
	}

	albumYamlMutex.Lock()
	defer albumYamlMutex.Unlock()

	// Write YAML data to file.
	if err := ioutil.WriteFile(fileName, data, os.ModePerm); err != nil {
		return err
	}

	return nil
}

// LoadFromYaml photo data from a YAML file.
func (m *Album) LoadFromYaml(fileName string) error {
	data, err := ioutil.ReadFile(fileName)

	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, m); err != nil {
		return err
	}

	return nil
}

// YamlFileName returns the YAML backup file name.
func (m *Album) YamlFileName(albumsPath string) string {
	return filepath.Join(albumsPath, m.AlbumType, m.AlbumUID+fs.YamlExt)
}
