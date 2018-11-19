package cmd

import (
	"github.com/nupplaphil/kopano-ldap/lib/kopano"
	"github.com/nupplaphil/kopano-ldap/lib/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"log"
	"os"
)

// userdeleteCmd represents the userdelete command
var userdeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runUserDelete(cmd.Flags())
	},
}

func init() {
	UserCmd.AddCommand(userdeleteCmd)

	userdeleteCmd.Flags().StringP("user", "u", "", "The username")
	userdeleteCmd.MarkFlagRequired("user")
}

func runUserDelete(flags *pflag.FlagSet) {
	host, port, fqdn, user, password := LdapFlags()
	baseDn := utils.GetBaseDN(fqdn)

	conn := kopano.Connect(host, port, fqdn, user, password)

	user, err := flags.GetString("user")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	kopano.Del(conn, baseDn, user)
}
