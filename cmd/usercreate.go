package cmd

import (
	"github.com/nupplaphil/kopano-ldap/kopano"
	"github.com/spf13/cobra"
)

// userCreateCmd represents the usercreate command
var userCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creating an user for Kopano in LDAP.",
	Long:  `Creating an user for Kopano in LDAP.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUserCreate(cmd)
	},
}

func init() {
	UserCmd.AddCommand(userCreateCmd)

	userCreateCmd.Flags().StringP("user", "u", "", "The name of the user. With this name the user will log on to the store.")
	userCreateCmd.Flags().StringP("fullname", "", "", "The full name of the user.")
	userCreateCmd.Flags().StringArray("email", nil, "The email address of the user. Often this is '<user name>@<email domain>'. You can define more than one email address, which will be set as an alias for this user")
	userCreateCmd.Flags().StringP("password", "p", "", "The password in plain text. The password will be stored encrypted in LDAP.")
	userCreateCmd.Flags().BoolP("active", "a", true, "The active state of this user. If set to 'yes' (or not set), the user is able to login, otherwise 'no'")
	userCreateCmd.MarkFlagRequired("user")
	userCreateCmd.MarkFlagRequired("fullname")
	userCreateCmd.MarkFlagRequired("email")
	userCreateCmd.MarkFlagRequired("password")
}

func runUserCreate(cmd *cobra.Command) error {
	flags := cmd.Flags()

	user, _ := flags.GetString("user")
	fullName, _ := flags.GetString("fullname")
	email, _ := flags.GetStringArray("email")
	password, _ := flags.GetString("password")
	active, _ := flags.GetBool("active")

	userSettings := kopano.NewUserSettings(user)
	userSettings.Fullname = fullName
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
		return err
	}

	if err := kopano.AddUser(client, baseDn, userSettings); err != nil {
		return err
	} else {
		cmd.Printf("user %q successfully created\n", user)
	}

	return nil
}
