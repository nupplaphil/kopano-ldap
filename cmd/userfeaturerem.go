package cmd

import (
	"github.com/nupplaphil/kopano-ldap/kopano"
	"github.com/spf13/pflag"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// userfeatureremCmd represents the userfeaturerem command
var userfeatureremCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removing a feature of an user in Kopano.",
	Long:  `Removing a feature of an user in Kopano.`,
	Run: func(cmd *cobra.Command, args []string) {
		runUserFeatureRemove(cmd.Flags())
	},
	TraverseChildren: true,
}

func init() {
	userfeatureCmd.AddCommand(userfeatureremCmd)

	userfeatureremCmd.Flags().StringArrayP("feature", "a", nil, "Removing one or more features (imap, pop3 or mobile)")
	userfeatureaddCmd.MarkFlagRequired("feature")

	userfeatureremCmd.Flags().StringP("user", "u", "", "The name of the user. With this name the user will log on to the store.")
	userfeatureaddCmd.MarkFlagRequired("user")
}

func runUserFeatureRemove(flags *pflag.FlagSet) {
	user, err := flags.GetString("user")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	features, err := flags.GetStringArray("feature")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW := LdapFlags()
	baseDn := kopano.GetBaseDN(ldapDomain)

	client, err := kopano.Connect(ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	kopano.RemoveUserFeatures(client, baseDn, user, features)
}
