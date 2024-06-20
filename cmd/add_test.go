package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestAddCommand(t *testing.T) {
	binaryPath, _, cleanup := SetupTest(t)
	defer cleanup()

	tests := []struct {
		name           string
		args           []string
		input          string
		expectedOutput string
		checkDB        bool
	}{
		{
			name:           "Add term successfully",
			args:           []string{"add", "-t", "testAddTerm", "testAddCategory"},
			input:          "\n",
			expectedOutput: "Term added successfully",
			checkDB:        true,
		},
		{
			name:           "Conflicting flags",
			args:           []string{"add", "-t", "testAddTerm", "testAddCategory", "-d"},
			expectedOutput: "Error: The -t and -d flags cannot be used together",
			checkDB:        false,
		},
		{
			name:           "Empty term name with -t flag",
			args:           []string{"add", "-t", "", "testAddCategory"},
			expectedOutput: "Error: The -t flag requires a non-empty term name",
			checkDB:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)

			if tt.input != "" {
				var stdin bytes.Buffer
				stdin.WriteString(tt.input)
				cmd.Stdin = &stdin
			}

			output, err := cmd.CombinedOutput()

			assert.Contains(t, string(output), tt.expectedOutput)

			if tt.checkDB {
				assert.NoError(t, err)
				terms, err := core.GetStorage().LoadAllData()
				assert.NoError(t, err)

				found := false
				for _, term := range terms {
					fmt.Printf("term name: %s, term category: %s\n", term.Name, term.Category)
					if term.Name == tt.args[2] && term.Category == strings.ToUpper(tt.args[3]) {
						found = true
						assert.True(t, term.Active)
						break
					}
				}
				assert.True(t, found, "Term not found in the database")
			} else {
				assert.Error(t, err, "Expected an error but got none")
			}
		})
	}
}
