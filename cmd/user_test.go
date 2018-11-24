// +build integration

package cmd

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestListAllUser(t *testing.T) {
	tests := []struct {
		args   []string
		output string
	}{
		{[]string{"user", "philipp7"}, `bla`},
	}

	for _, test := range tests {
		output, err := executeCommand(rootCmd, test.args...)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		assert.Equal(t, output, test.output)
	}
}
