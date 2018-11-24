package cmd

import (
	"github.com/nupplaphil/kopano-ldap/kopano"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// UserCmd represents the user command
var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "Creating, modifying or creating users.",
	Long:  `Creating, modifying or creating users.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUser(cmd.Flags(), args)
	},
}

func init() {
	rootCmd.AddCommand(UserCmd)

	UserCmd.Flags().BoolP("list", "l", false, "List all users")
}

func runUser(flags *pflag.FlagSet, args []string) error {
	list, _ := flags.GetBool("list")
	if list {
		return listAllUser()
	}

	argsLen := len(args)

	if argsLen == 1 {
		return listUser(args[0])
	}

	return listAllUser()
}

func listAllUser() error {
	ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW := LdapFlags()
	base := kopano.GetBaseDN(ldapDomain)

	client, err := kopano.Connect(ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW)
	if err != nil {
		return err
	}

	if err = kopano.ListAll(client, base); err != nil {
		return err
	}

	return nil
}

func listUser(user string) error {
	ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW := LdapFlags()
	baseDn := kopano.GetBaseDN(ldapDomain)

	client, err := kopano.Connect(ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW)
	if err != nil {
		return err
	}
	if err = kopano.ListUser(client, baseDn, user); err != nil {
		return err
	}

	return nil
}
