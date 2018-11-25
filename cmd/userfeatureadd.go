package cmd

import (
	"github.com/nupplaphil/kopano-ldap/kopano"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// userfeatureaddCmd represents the userfeatureadd command
var userfeatureaddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adding a new feature to an user in Kopano.",
	Long:  `Adding a new feature to an user in Kopano.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUserFeatureAdd(cmd.Flags())
	},
}

func init() {
	userfeatureCmd.AddCommand(userfeatureaddCmd)

	userfeatureaddCmd.Flags().StringArrayP("feature", "a", nil, "Adding one or more features (imap, pop3 or mobile)")
	userfeatureaddCmd.MarkFlagRequired("feature")

	userfeatureaddCmd.Flags().StringP("user", "u", "", "The name of the user. With this name the user will log on to the store.")
	userfeatureaddCmd.MarkFlagRequired("user")
}

func runUserFeatureAdd(flags *pflag.FlagSet) error {
	user, _ := flags.GetString("user")
	features, _ := flags.GetStringArray("feature")

	ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW := LdapFlags()
	baseDn := kopano.GetBaseDN(ldapDomain)

	client, err := kopano.Connect(ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW)
	if err != nil {
		return err
	}

	if err := kopano.AddUserFeatures(client, baseDn, user, features); err != nil {
		return err
	}

	return nil
}
