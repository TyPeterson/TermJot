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

func TestDoneCommand(t *testing.T) {
	binaryPath, _, cleanup := SetupTest(t)
	defer cleanup()

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
			name: "Mark existing term as done successfully",
			setup: func() {
				addTerm("existingTerm", "testCategory")
			},
			args:           []string{"done", "testCategory", "-t", "existingTerm"},
			expectedOutput: "Term marked as done",
			checkDB:        true,
		},
		{
			name: "Mark non-existing term as done",
			setup: func() {
				addTerm("existingTerm", "testCategory")
			},
			args:             []string{"done", "testCategory", "-t", "nonExistingTerm"},
			expectedOutput:   "Error: Term not found",
			unexpectedOutput: "Term marked as done",
			checkDB:          false,
		},
		{
			name:             "Mark term as done in non-existing category",
			args:             []string{"done", "nonExistingCategory", "-t", "someTerm"},
			expectedOutput:   "Error: Category not found",
			unexpectedOutput: "Term marked as done",
			checkDB:          false,
		},
		{
			name:           "Mark as done with -t flag but no term",
			args:           []string{"done", "testCategory", "-t", ""},
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
						assert.False(t, term.Active, "Term should be marked as done")
						break
					}
				}
				assert.True(t, found, "Term should be found in the database")
			}

			if tt.expectedOutput == "Error: The -t flag requires a non-empty term name" {
				assert.Error(t, err, "Expected an error but got none")
			} else if strings.Contains(tt.expectedOutput, "Error:") {
				assert.Contains(t, output, tt.expectedOutput)
			} else {
				assert.NoError(t, err, "Unexpected error occurred")
			}
		})
	}
}
