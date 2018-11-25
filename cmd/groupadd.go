package cmd

import (
	"github.com/nupplaphil/kopano-ldap/kopano"
	"github.com/spf13/cobra"
)

// userAddCmd represents the groupadd command
var groupAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adding an user to a group for Kopano in LDAP.",
	Long:  `Adding an user to a group for Kopano in LDAP.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGroupAdd(cmd)
	},
}

func init() {
	groupCmd.AddCommand(groupAddCmd)

	groupAddCmd.Flags().StringP("name", "n", "", "The name of the group.")
	groupAddCmd.Flags().StringArrayP("user", "u", []string{}, "The user to add.")
	groupAddCmd.MarkFlagRequired("name")
	groupAddCmd.MarkFlagRequired("user")
}

func runGroupAdd(cmd *cobra.Command) error {
	flags := cmd.Flags()

	group, _ := flags.GetString("group")
	userList, _ := flags.GetStringArray("user")

	ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW := LdapFlags()
	baseDn := kopano.GetBaseDN(ldapDomain)

	client, err := kopano.Connect(ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW)
	if err != nil {
		return err
	}

	if err := kopano.AddUserToGroup(client, baseDn, group, userList); err != nil {
		return err
	} else {
		cmd.Printf("users %v successfully added to group %q\n", userList, group)
	}

	return nil
}
