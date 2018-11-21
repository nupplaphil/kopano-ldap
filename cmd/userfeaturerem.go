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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runUserFeatureRemove(cmd.Flags())
	},
	TraverseChildren: true,
}

func init() {
	userfeatureCmd.AddCommand(userfeatureremCmd)

	userfeatureremCmd.Flags().StringP("user", "u", "", "The user name of the user")
	userfeatureaddCmd.MarkFlagRequired("user")

	userfeatureremCmd.Flags().StringArrayP("feature", "a", nil, "Adding features")
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
