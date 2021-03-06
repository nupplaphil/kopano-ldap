package cmd

import (
	"github.com/nupplaphil/kopano-ldap/kopano"
	"github.com/spf13/cobra"
	"io"
)

// UserCmd represents the user command
var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "Creating, modifying or creating users.",
	Long:  `Creating, modifying or creating users.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUser(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(UserCmd)

	UserCmd.Flags().BoolP("list", "l", false, "List all users")
}

func runUser(cmd *cobra.Command, args []string) error {
	list, _ := cmd.Flags().GetBool("list")
	if list {
		return listAllUser(cmd.OutOrStdout())
	}

	argsLen := len(args)

	if argsLen == 1 {
		return listUser(args[0], cmd.OutOrStdout())
	}

	return listAllUser(cmd.OutOrStdout())
}

func listAllUser(writer io.Writer) error {
	ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW := LdapFlags()
	base := kopano.GetBaseDN(ldapDomain)

	client, err := kopano.Connect(ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW)
	if err != nil {
		return err
	}

	if err = kopano.ListAllUsers(client, base, writer); err != nil {
		return err
	}

	return nil
}

func listUser(user string, writer io.Writer) error {
	ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW := LdapFlags()
	baseDn := kopano.GetBaseDN(ldapDomain)

	client, err := kopano.Connect(ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW)
	if err != nil {
		return err
	}
	if err = kopano.ListUser(client, baseDn, user, writer); err != nil {
		return err
	}

	return nil
}
