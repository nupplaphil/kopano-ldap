package cmd

import (
	"github.com/nupplaphil/kopano-ldap/kopano"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"log"
	"os"
)

// userfeatureaddCmd represents the userfeatureadd command
var userfeatureaddCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runUserFeatureAdd(cmd.Flags())
	},
}

func init() {
	userfeatureCmd.AddCommand(userfeatureaddCmd)

	userfeatureaddCmd.Flags().StringArrayP("feature", "a", nil, "Adding features")

	userfeatureaddCmd.Flags().StringP("user", "u", "", "The user name of the user")
	userfeatureaddCmd.MarkFlagRequired("user")
}

func runUserFeatureAdd(flags *pflag.FlagSet) {
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

	kopano.AddUserFeatures(client, baseDn, user, features)
}
