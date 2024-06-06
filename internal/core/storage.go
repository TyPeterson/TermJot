package core

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/TyPeterson/TermJot/models"
)

const dataFileName = "termjot_data.json"

type Storage struct {
	FilePath string
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

// ------------- LoadData() -------------
func (s *Storage) LoadData() ([]models.Term, []models.Category, error) {
	if _, err := os.Stat(s.FilePath); os.IsNotExist(err) {
		return []models.Term{}, []models.Category{}, nil
	}

	data, err := ioutil.ReadFile(s.FilePath)
	if err != nil {
		return nil, nil, err
	}

	var terms []models.Term
	var categories []models.Category
	if err := json.Unmarshal(data, &struct {
		Terms      *[]models.Term     `json:"terms"`
		Categories *[]models.Category `json:"categories"`
	}{&terms, &categories}); err != nil {
		return nil, nil, err
	}

	return terms, categories, nil
}

// ------------- SaveData() -------------
func (s *Storage) SaveData(terms []models.Term, categories []models.Category) error {
	data, err := json.MarshalIndent(struct {
		Terms      []models.Term     `json:"terms"`
		Categories []models.Category `json:"categories"`
	}{terms, categories}, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(s.FilePath, data, 0644)
}
