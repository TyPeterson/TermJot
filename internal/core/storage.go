package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"github.com/TyPeterson/TermJot/models"
)

const dataFileName = "termjot_data.json"

type Storage struct {
	FilePath string
}


// ------------- Init() -------------
func Init() error {
	storage, err := NewStorage()

    if err != nil {
        return err
    }

	loadedTerms := storage.LoadData()
	terms = loadedTerms

    return nil
}


// ------------- NewStorage() -------------
func NewStorage() (*Storage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	filePath := filepath.Join(homeDir, dataFileName)
	return &Storage{FilePath: filePath}, nil
}

// ------------- Save -------------
func Save() error {
    storage, err := NewStorage()
    if err != nil {
        return err
    }
    return storage.SaveData(terms)
}


// ------------- LoadData() -------------
func (s *Storage) LoadData() []models.Term {
	if _, err := os.Stat(s.FilePath); os.IsNotExist(err) {
		return []models.Term{}
	}

	data, err := os.ReadFile(s.FilePath)
	if err != nil {
		return nil
	}

	var terms []models.Term
	if err := json.Unmarshal(data, &struct {
		Terms      *[]models.Term     `json:"terms"`
	}{&terms}); err != nil {
		return nil
	}

	return terms
}

// ------------- SaveData() -------------
func (s *Storage) SaveData(terms []models.Term) error {
	data, err := json.MarshalIndent(struct {
		Terms      []models.Term     `json:"terms"`
	}{terms}, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.FilePath, data, 0644)
}
