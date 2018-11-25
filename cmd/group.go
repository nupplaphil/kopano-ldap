package cmd

import (
	"github.com/nupplaphil/kopano-ldap/kopano"
	"github.com/spf13/cobra"
	"io"
)

// groupCmd represents the group command
var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Creating, modifying or creating groups.",
	Long:  `Creating, modifying or creating groups.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGroup(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(groupCmd)

	groupCmd.Flags().BoolP("list", "l", false, "List all groups")
}

func runGroup(cmd *cobra.Command, args []string) error {
	list, _ := cmd.Flags().GetBool("list")
	if list {
		return listAllGroups(cmd.OutOrStdout())
	}

	argsLen := len(args)

	if argsLen == 1 {
		return listGroup(args[0], cmd.OutOrStdout())
	}

	return listAllGroups(cmd.OutOrStdout())
}

func listAllGroups(writer io.Writer) error {
	ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW := LdapFlags()
	base := kopano.GetBaseDN(ldapDomain)

	client, err := kopano.Connect(ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW)
	if err != nil {
		return err
	}

	if err = kopano.ListAllGroups(client, base, writer); err != nil {
		return err
	}

	return nil
}

func listGroup(group string, writer io.Writer) error {
	ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW := LdapFlags()
	baseDn := kopano.GetBaseDN(ldapDomain)

	client, err := kopano.Connect(ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW)
	if err != nil {
		return err
	}
	if err = kopano.ListGroup(client, baseDn, group, writer); err != nil {
		return err
	}

	return nil
}
