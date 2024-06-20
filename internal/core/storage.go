package core

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const dbFileName = ".termjot.db"

type Storage struct {
	DB *sql.DB
}

var storage *Storage

// ------------- Init -------------
func Init() error {
	dbPath := os.Getenv("TERMJOT_DB_PATH")
	if dbPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		dbPath = filepath.Join(homeDir, dbFileName)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS terms (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		definition TEXT,
		category TEXT,
		active BOOLEAN NOT NULL CHECK (active IN (0, 1))
	)`)
	if err != nil {
		return err
	}

	storage = &Storage{DB: db}
	return nil
}

// ------------- SaveData -------------
func (s *Storage) SaveData(term Term) error {
	_, err := s.DB.Exec("INSERT OR REPLACE INTO terms (name, definition, category, active) VALUES (?, ?, ?, ?)",
		term.Name, term.Definition, term.Category, term.Active)
	return err
}

// ------------- LoadAllData -------------
func (s *Storage) LoadAllData() ([]Term, error) {
	rows, err := s.DB.Query("SELECT name, definition, category, active FROM terms")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var terms []Term
	for rows.Next() {
		var term Term
		if err := rows.Scan(&term.Name, &term.Definition, &term.Category, &term.Active); err != nil {
			return nil, err
		}
		terms = append(terms, term)
	}
	return terms, nil
}

// ------------- RemoveData -------------
func (s *Storage) RemoveData(term Term) error {
	_, err := s.DB.Exec("DELETE FROM terms WHERE name = ? AND category = ?", term.Name, term.Category)
	return err
}

// ------------- UpdateData -------------
func (s *Storage) UpdateData(term Term) error {
	_, err := s.DB.Exec("UPDATE terms SET definition = ?, category = ?, active = ? WHERE name = ?",
		term.Definition, term.Category, term.Active, term.Name)
	return err
}

// ------------- Close -------------
func (s *Storage) Close() error {
	return s.DB.Close()
}

// ------------- SetStorage -------------
func SetStorage(s *Storage) {
	storage = s
}

// ------------- GetStorage -------------
func GetStorage() *Storage {
	return storage
}
