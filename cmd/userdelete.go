package cmd

import (
	"github.com/nupplaphil/kopano-ldap/kopano"
	"github.com/spf13/cobra"
)

// userdeleteCmd represents the userdelete command
var userdeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deleting an user for Kopano in LDAP.",
	Long:  `Creating an user for Kopano in LDAP.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUserDelete(cmd)
	},
}

func init() {
	UserCmd.AddCommand(userdeleteCmd)

	userdeleteCmd.Flags().StringP("user", "u", "", "The name of the user.")
	userdeleteCmd.MarkFlagRequired("user")
}

func runUserDelete(cmd *cobra.Command) error {
	ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW := LdapFlags()
	baseDn := kopano.GetBaseDN(ldapDomain)

	client, err := kopano.Connect(ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW)
	if err != nil {
		return err
	}

	user, err := cmd.Flags().GetString("user")
	if err != nil {
		return err
	}

	if err := kopano.Del(client, baseDn, user); err != nil {
		return err
	} else {
		cmd.Printf("user %q successfully deleted.", user)
	}

	return nil
}
