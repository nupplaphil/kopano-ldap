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
	userCmd.AddCommand(userdeleteCmd)

	userdeleteCmd.Flags().StringP("user", "u", "", "The username")
	userdeleteCmd.MarkFlagRequired("user")
}

func runUserDelete(flags *pflag.FlagSet) {
	host, port, fqdn, user, password := ldapFlags()
	baseDn := utils.GetBaseDN(fqdn)

	conn := ldap.Connect(host, port, fqdn, user, password)

	user, err := flags.GetString("user")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	kopano.Del(conn, baseDn, user)
}
