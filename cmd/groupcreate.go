package cmd

import (
	"github.com/nupplaphil/kopano-ldap/kopano"
	"github.com/spf13/cobra"
)

// userCreateCmd represents the usercreate command
var groupCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creating an user for Kopano in LDAP.",
	Long:  `Creating an user for Kopano in LDAP.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGroupCreate(cmd)
	},
}

func init() {
	groupCmd.AddCommand(groupCreateCmd)

	groupCreateCmd.Flags().StringP("name", "n", "", "The name of the group.")
	groupCreateCmd.Flags().BoolP("active", "a", true, "The active state of this group.")
	groupCreateCmd.Flags().BoolP("security", "", false, "The active state of this group.")
	groupCreateCmd.Flags().BoolP("hidden", "", false, "The active state of this group.")
	groupCreateCmd.MarkFlagRequired("name")
}

func runGroupCreate(cmd *cobra.Command) error {
	flags := cmd.Flags()

	name, _ := flags.GetString("name")
	active, _ := flags.GetBool("active")
	security, _ := flags.GetBool("security")
	hidden, _ := flags.GetBool("hidden")

	groupSettings := kopano.NewGroupSettings(name)
	groupSettings.Active = active
	groupSettings.Security = security
	groupSettings.Hidden = hidden

	ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW := LdapFlags()
	baseDn := kopano.GetBaseDN(ldapDomain)

	client, err := kopano.Connect(ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW)
	if err != nil {
		return err
	}

	if err := kopano.AddGroup(client, baseDn, groupSettings); err != nil {
		return err
	} else {
		cmd.Printf("group %q successfully created\n", name)
	}

	return nil
}
