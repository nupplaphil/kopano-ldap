package cmd

import (
	"github.com/nupplaphil/kopano-ldap/kopano"
	"github.com/spf13/cobra"
)

// groupDeleteCmd represents the groupdelete command
var groupDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deleting a group for Kopano in LDAP.",
	Long:  `Creating a group for Kopano in LDAP.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGroupDelete(cmd)
	},
}

func init() {
	UserCmd.AddCommand(groupDeleteCmd)

	groupDeleteCmd.Flags().StringP("name", "n", "", "The name of the group.")
	groupDeleteCmd.MarkFlagRequired("name")
}

func runGroupDelete(cmd *cobra.Command) error {
	ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW := LdapFlags()
	baseDn := kopano.GetBaseDN(ldapDomain)

	client, err := kopano.Connect(ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW)
	if err != nil {
		return err
	}

	group, _ := cmd.Flags().GetString("name")

	if err := kopano.DelGroup(client, baseDn, group); err != nil {
		return err
	} else {
		cmd.Printf("group %q successfully deleted.", group)
	}

	return nil
}
