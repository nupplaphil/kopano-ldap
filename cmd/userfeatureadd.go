// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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

	userfeatureaddCmd.Flags().StringP("user", "u", "", "The user name of the user")

	userfeatureaddCmd.Flags().StringArrayP("feature", "a", nil, "Adding features")
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

	host, port, fqdn, ldap_user, password := ldapFlags()
	baseDn := utils.GetBaseDN(fqdn)

	conn := ldap.Connect(host, port, fqdn, ldap_user, password)

	kopano.AddUserFeatures(conn, baseDn, user, features)
}
