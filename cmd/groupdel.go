package cmd

import (
	"github.com/nupplaphil/kopano-ldap/kopano"
	"github.com/spf13/cobra"
)

// userAddCmd represents the groupadd command
var groupDelCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deleting an user from a group for Kopano in LDAP.",
	Long:  `Deleting an user from a group for Kopano in LDAP.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGroupDel(cmd)
	},
}

func init() {
	groupCmd.AddCommand(groupDelCmd)

	groupDelCmd.Flags().StringP("name", "n", "", "The name of the group.")
	groupDelCmd.Flags().StringArrayP("user", "u", []string{}, "The user to add.")
	groupDelCmd.MarkFlagRequired("name")
	groupDelCmd.MarkFlagRequired("user")
}

func runGroupDel(cmd *cobra.Command) error {
	flags := cmd.Flags()

	group, _ := flags.GetString("group")
	userList, _ := flags.GetStringArray("user")

	ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW := LdapFlags()
	baseDn := kopano.GetBaseDN(ldapDomain)

	client, err := kopano.Connect(ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW)
	if err != nil {
		return err
	}

	if err := kopano.DelUserFromGroup(client, baseDn, group, userList); err != nil {
		return err
	} else {
		cmd.Printf("users %v successfully deleted from group %q\n", userList, group)
	}

	return nil
}
