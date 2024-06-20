package cmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestRemoveCommand(t *testing.T) {
	binaryPath, _, cleanup := SetupTest(t)
	defer cleanup()

	// Helper function to add a term for testing removal
	addTerm := func(termName, category string) {
		core.HandleAdd(termName, category)
	}

	tests := []struct {
		name             string
		setup            func()
		args             []string
		input            string
		expectedOutput   string
		unexpectedOutput string
		checkDB          bool
	}{
		{
			name: "Remove existing term successfully",
			setup: func() {
				addTerm("existingTerm", "testCategory")
			},
			args:           []string{"remove", "testCategory", "-t", "existingTerm"},
			expectedOutput: "Term removed successfully",
			checkDB:        true,
		},
		{
			name: "Remove non-existing term",
			setup: func() {
				addTerm("existingTerm", "testCategory")
			},
			args:             []string{"remove", "testCategory", "-t", "nonExistingTerm"},
			expectedOutput:   "Error: Term not found",
			unexpectedOutput: "Term removed successfully",
			checkDB:          false,
		},
		{
			name:             "Remove from non-existing category",
			args:             []string{"remove", "nonExistingCategory", "-t", "someTerm"},
			expectedOutput:   "Error: Category not found",
			unexpectedOutput: "Term removed successfully",
			checkDB:          false,
		},
		{
			name:           "Remove with -t flag but no term",
			args:           []string{"remove", "testCategory", "-t", ""},
			expectedOutput: "Error: The -t flag requires a non-empty term name",
			checkDB:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			cmd := exec.Command(binaryPath, tt.args...)

			// Capture both stdout and stderr
			var stdoutBuf, stderrBuf bytes.Buffer
			cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
			cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

			if tt.input != "" {
				cmd.Stdin = strings.NewReader(tt.input)
			}

			err := cmd.Run()
			output := stdoutBuf.String() + stderrBuf.String()

			if tt.expectedOutput != "" {
				assert.Contains(t, output, tt.expectedOutput)
			}
			if tt.unexpectedOutput != "" {
				assert.NotContains(t, output, tt.unexpectedOutput)
			}

			if tt.checkDB {
				terms, err := core.GetStorage().LoadAllData()
				assert.NoError(t, err)

				found := false
				for _, term := range terms {
					if term.Name == tt.args[3] && term.Category == strings.ToUpper(tt.args[1]) {
						found = true
						break
					}
				}
				assert.False(t, found, "Term should not be found in the database after removal")
			}

			// Check for error only if we're expecting the command itself to fail
			if tt.expectedOutput == "Error: The -t flag requires a non-empty term name" {
				assert.Error(t, err, "Expected an error but got none")
				assert.Contains(t, output, tt.expectedOutput)
			} else if strings.Contains(tt.expectedOutput, "Error:") {
				assert.Contains(t, output, tt.expectedOutput)
			} else {
				assert.NoError(t, err, "Unexpected error occurred")
				assert.Contains(t, output, tt.expectedOutput)
			}
		})
	}
}
