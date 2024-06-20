package cmd

import (
	"database/sql"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/TyPeterson/TermJot/internal/core"
)

func SetupTest(t *testing.T) (string, string, func()) {
	binaryPath := BuildBinary(t)
	db, dbPath, err := InitTestDB()
	if err != nil {
		t.Fatalf("Failed to initialize test db: %v", err)
	}

	os.Setenv("TERMJOT_DB_PATH", dbPath)

	err = core.Init()
	if err != nil {
		t.Fatalf("Failed to initialize core: %v", err)
	}

	cleanup := func() {
		os.Remove(binaryPath)
		os.Remove(dbPath)
		db.Close()
		os.Unsetenv("TERMJOT_DB_PATH")
	}

	return binaryPath, dbPath, cleanup
}

func BuildBinary(t *testing.T) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	binaryPath := filepath.Join(wd, "jot")
	cmd := exec.Command("go", "build", "-o", binaryPath, "../main.go")
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	return binaryPath
}

func InitTestDB() (*core.Storage, string, error) {
	tmpfile, err := os.CreateTemp("", "termjot-test-*.db")
	if err != nil {
		return nil, "", err
	}

	db, err := sql.Open("sqlite3", tmpfile.Name())
	if err != nil {
		return nil, "", err
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
		return nil, "", err
	}

	storage := &core.Storage{DB: db}
	core.SetStorage(storage)
	return storage, tmpfile.Name(), nil
}
