// Copyright © 2018 Philipp Holzer <admin@philipp.info>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/nupplaphil/kopano-ldap/ldap"
	"github.com/nupplaphil/kopano-ldap/ldap/kopano"
	"github.com/nupplaphil/kopano-ldap/ldap/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runUser(cmd.Flags(), args)
	},
}

func init() {
	rootCmd.AddCommand(userCmd)

	userCmd.Flags().BoolP("list", "l", false, "ListAll all users")
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
	host, port, fqdn, user, password := ldapFlags()
	base := utils.GetBaseDN(fqdn)

	conn := ldap.Connect(host, port, fqdn, user, password)
	kopano.ListAll(conn, base)
}

func listUser(user string) {
	host, port, fqdn, ldap_user, password := ldapFlags()
	baseDn := utils.GetBaseDN(fqdn)

	conn := ldap.Connect(host, port, fqdn, ldap_user, password)
	kopano.ListUser(conn, baseDn, user)
}
