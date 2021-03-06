package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var versionStr string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kopano-ld",
	Short: "This CLI is built for the LDAP User Management of Kopano.",
	Long: `The kopano-ld is an administration tool for managing user and groups in LDAP.

The tool can be used to get more information about users and groups too.`,
	Run: func(cmd *cobra.Command, args []string) {
		runRoot(cmd)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	versionStr = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kopano-kopano.yaml)")

	rootCmd.PersistentFlags().String("ldaphost", "localhost", "LDAP host")
	rootCmd.PersistentFlags().Int("ldapport", 389, "LDAP port")
	rootCmd.PersistentFlags().StringP("domain", "b", "example.org", "LDAP Base domain")
	rootCmd.PersistentFlags().String("ldapuser", "admin", "LDAP user")
	rootCmd.PersistentFlags().String("ldappass", "admin", "LDAP password")

	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("ldaphost"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("ldapport"))
	viper.BindPFlag("domain", rootCmd.PersistentFlags().Lookup("domain"))
	viper.BindPFlag("admin_user", rootCmd.PersistentFlags().Lookup("ldapuser"))
	viper.BindPFlag("admin_password", rootCmd.PersistentFlags().Lookup("ldappass"))

	rootCmd.Flags().BoolP("version", "v", false, "Printing out the current version")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".kopano-kopano" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".kopano-kopano")
	}

	viper.SetEnvPrefix("LDAP")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// LdapFlags returns all LDAP properties based on either parameters or environment variables
func LdapFlags() (string, int, string, string, string) {
	viper.SetEnvPrefix("LDAP")
	viper.AutomaticEnv()

	ldapHost := viper.GetString("host")
	ldapPort := viper.GetInt("port")
	ldapDomain := viper.GetString("domain")
	ldapUser := viper.GetString("admin_user")
	ldapPW := viper.GetString("admin_password")

	return ldapHost, ldapPort, ldapDomain, ldapUser, ldapPW
}

func runRoot(cmd *cobra.Command) {
	versionFlag, _ := cmd.Flags().GetBool("version")

	if versionFlag {
		cmd.Println(versionStr)
	} else {
		cmd.Help()
	}
}
