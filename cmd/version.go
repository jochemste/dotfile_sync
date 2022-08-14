package cmd

import (
	"fmt"

	"github.com/jochemste/dotfile_sync/libdotfilesync"

	"github.com/spf13/cobra"
)

var cmdVersion *cobra.Command

func runVersion(cmd *cobra.Command, args []string) error {
	verbose, err := rootCmd.Flags().GetBool("verbose")
	if err != nil {
		return err
	}

	if verbose == true {
		fmt.Println(libdotfilesync.NAME + " - v" + libdotfilesync.VERSION + " by " + libdotfilesync.AUTHOR)
	} else {
		fmt.Println(libdotfilesync.VERSION)
	}
	return nil
}

func init() {
	cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		RunE:  runVersion,
	}

	rootCmd.AddCommand(cmdVersion)
}
