package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kopano-ld",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kopano-lib.yaml)")

	rootCmd.PersistentFlags().String("ldaphost", "localhost", "LDAP host")
	rootCmd.PersistentFlags().Int("ldapport", 389, "LDAP port")
	rootCmd.PersistentFlags().StringP("domain", "b", "example.org", "LDAP Base domain")
	rootCmd.PersistentFlags().String("ldapuser", "admin", "LDAP user")
	rootCmd.PersistentFlags().String("ldappass", "", "LDAP password")

	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("ldaphost"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("ldapport"))
	viper.BindPFlag("domain", rootCmd.PersistentFlags().Lookup("domain"))
	viper.BindPFlag("admin_user", rootCmd.PersistentFlags().Lookup("ldapuser"))
	viper.BindPFlag("admin_password", rootCmd.PersistentFlags().Lookup("ldappass"))
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

		// Search config in home directory with name ".kopano-lib" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".kopano-lib")
	}

	viper.SetEnvPrefix("LDAP")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func LdapFlags() (string, int, string, string, string) {
	viper.SetEnvPrefix("LDAP")
	viper.AutomaticEnv()

	host := viper.GetString("host")
	port := viper.GetInt("port")
	fqdn := viper.GetString("domain")
	user := viper.GetString("admin_user")
	password := viper.GetString("admin_password")

	return host, port, fqdn, user, password
}
