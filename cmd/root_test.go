package cmd

import (
	"bytes"
	"github.com/magiconair/properties/assert"
	"github.com/spf13/cobra"
	"os"
	"testing"
)

func TestLdapFlagsEnvironment(t *testing.T) {
	os.Setenv("LDAP_HOST", "test.host")
	os.Setenv("LDAP_PORT", "123")
	os.Setenv("LDAP_DOMAIN", "example.test")
	os.Setenv("LDAP_ADMIN_USER", "test_user")
	os.Setenv("LDAP_ADMIN_PASSWORD", "test_pw")

	ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW := LdapFlags()

	assert.Equal(t, ldapHost, "test.host")
	assert.Equal(t, ldapPort, 123)
	assert.Equal(t, ldapDomain, "example.test")
	assert.Equal(t, ldapUser, "test_user")
	assert.Equal(t, ldapPW, "test_pw")
}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOutput(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func TestExecute(t *testing.T) {
	tests := []struct {
		args   []string
		output string
	}{
		{[]string{}, `The kopano-ld is an administration tool for managing user and groups in LDAP.

The tool can be used to get more information about users and groups too.

Usage:
  kopano-ld [flags]
  kopano-ld [command]

Available Commands:
  help        Help about any command
  user        Creating, modifying or creating users.

Flags:
      --config string     config file (default is $HOME/.kopano-kopano.yaml)
  -b, --domain string     LDAP Base domain (default "example.org")
  -h, --help              help for kopano-ld
      --ldaphost string   LDAP host (default "localhost")
      --ldappass string   LDAP password (default "admin")
      --ldapport int      LDAP port (default 389)
      --ldapuser string   LDAP user (default "admin")
  -v, --version           Printing out the current version

Use "kopano-ld [command] --help" for more information about a command.
`},
	}

	for _, test := range tests {
		output, err := executeCommand(rootCmd, test.args...)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		assert.Equal(t, output, test.output)
	}
}
