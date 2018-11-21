package cmd

import (
	"github.com/nupplaphil/kopano-ldap/kopano"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"log"
	"os"
)

// userdeleteCmd represents the userdelete command
var userdeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deleting an user for Kopano in LDAP.",
	Long:  `Creating an user for Kopano in LDAP.`,
	Run: func(cmd *cobra.Command, args []string) {
		runUserDelete(cmd.Flags())
	},
}

func init() {
	UserCmd.AddCommand(userdeleteCmd)

	userdeleteCmd.Flags().StringP("user", "u", "", "The name of the user.")
	userdeleteCmd.MarkFlagRequired("user")
}

func runUserDelete(flags *pflag.FlagSet) {
	ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW := LdapFlags()
	baseDn := kopano.GetBaseDN(ldapDomain)

	client, err := kopano.Connect(ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	user, err := flags.GetString("user")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	kopano.Del(client, baseDn, user)
}
