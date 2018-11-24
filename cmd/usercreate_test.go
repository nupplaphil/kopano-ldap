// +build integration

package cmd

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestRunUserCreate(t *testing.T) {
	tests := []struct {
		args   []string
		output string
	}{
		{[]string{"user",
			"create",
			"--user",
			"johndoe",
			"--fullname",
			"John Doe",
			"--email",
			"john@doe.com",
			"--email",
			"john2@doe.com",
			"--password",
			"testpw",
		}, `bla`},
	}

	for _, test := range tests {
		output, err := executeCommand(rootCmd, test.args...)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		assert.Equal(t, output, test.output)
	}
}
