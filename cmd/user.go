package cmd

import (
	"github.com/nupplaphil/kopano-ldap/kopano"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"log"
	"os"
)

// UserCmd represents the user command
var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "Creating, modifying or creating users.",
	Long:  `Creating, modifying or creating users.`,
	Run: func(cmd *cobra.Command, args []string) {
		runUser(cmd.Flags(), args)
	},
}

func init() {
	rootCmd.AddCommand(UserCmd)

	UserCmd.Flags().BoolP("list", "l", false, "List all users")
}

func runUser(flags *pflag.FlagSet, args []string) {
	list, _ := flags.GetBool("list")
	if list {
		listAllUser()
		os.Exit(0)
	}

	argsLen := len(args)

	if argsLen == 1 {
		listUser(args[0])
		os.Exit(0)
	}

	listAllUser()
}

func listAllUser() {
	ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW := LdapFlags()
	base := kopano.GetBaseDN(ldapDomain)

	client, err := kopano.Connect(ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = kopano.ListAll(client, base)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func listUser(user string) {
	ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW := LdapFlags()
	baseDn := kopano.GetBaseDN(ldapDomain)

	client, err := kopano.Connect(ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	err = kopano.ListUser(client, baseDn, user)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
