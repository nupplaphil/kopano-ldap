package cmd

import (
	"github.com/nupplaphil/kopano-ldap/kopano"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"log"
	"os"
)

// usercreateCmd represents the usercreate command
var usercreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creating an user for Kopano in LDAP.",
	Long:  `Creating an user for Kopano in LDAP.`,
	Run: func(cmd *cobra.Command, args []string) {
		runUserCreate(cmd.Flags())
	},
}

func init() {
	UserCmd.AddCommand(usercreateCmd)

	usercreateCmd.Flags().StringP("user", "u", "", "The name of the user. With this name the user will log on to the store.")
	usercreateCmd.Flags().StringP("fullname", "", "", "The full name of the user.")
	usercreateCmd.Flags().StringArray("email", nil, "The email address of the user. Often this is '<user name>@<email domain>'. You can define more than one email address, which will be set as an alias for this user")
	usercreateCmd.Flags().StringP("password", "p", "", "The password in plain text. The password will be stored encrypted in LDAP.")
	usercreateCmd.Flags().BoolP("active", "a", true, "The active state of this user. If set to 'yes' (or not set), the user is able to login, otherwise 'no'")
	usercreateCmd.MarkFlagRequired("user")
	usercreateCmd.MarkFlagRequired("fullname")
	usercreateCmd.MarkFlagRequired("email")
	usercreateCmd.MarkFlagRequired("password")
}

func runUserCreate(flags *pflag.FlagSet) {
	user, err := flags.GetString("user")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fullname, err := flags.GetString("fullname")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	email, err := flags.GetStringArray("email")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	password, err := flags.GetString("password")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	active, err := flags.GetBool("active")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	userSettings := kopano.NewUserSettings(user)
	userSettings.Fullname = fullname
	userSettings.Email = email[0]
	userSettings.Password = password
	userSettings.Active = active

	if len(email) > 1 {
		userSettings.Aliases = email[1:]
	}

	ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW := LdapFlags()
	baseDn := kopano.GetBaseDN(ldapDomain)

	client, err := kopano.Connect(ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	kopano.Add(client, baseDn, userSettings)
}
