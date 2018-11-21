package cmd

import (
	"github.com/spf13/cobra"
)

// userfeatureCmd represents the userfeature command
var userfeatureCmd = &cobra.Command{
	Use:   "feature",
	Short: "Modifying features of an user in Kopano.",
	Long:  `Modifying features of an user in Kopano.`,
}

func init() {
	UserCmd.AddCommand(userfeatureCmd)
}
