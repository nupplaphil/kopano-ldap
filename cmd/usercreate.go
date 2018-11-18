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
	userCmd.AddCommand(usercreateCmd)

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

	host, port, fqdn, user, password := ldapFlags()
	baseDn := utils.GetBaseDN(fqdn)

	conn := ldap.Connect(host, port, fqdn, user, password)

	kopano.Add(conn, baseDn, userSettings)
}
