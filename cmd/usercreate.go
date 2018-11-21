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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runUserCreate(cmd.Flags())
	},
}

func init() {
	UserCmd.AddCommand(usercreateCmd)

	usercreateCmd.Flags().StringP("user", "u", "", "The user name of the user")
	usercreateCmd.Flags().StringP("fullname", "", "", "The full name of the user")
	usercreateCmd.Flags().StringArray("email", nil, "The email address of the user (The first one is the main address)")
	usercreateCmd.Flags().StringP("password", "p", "", "The password of the user")
	usercreateCmd.Flags().BoolP("active", "a", true, "The active state of the user")
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
		userSettings.Aliase = email[1:]
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
