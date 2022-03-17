package cmd

import (
	//"log"
	"fmt"

	"github.com/jochemste/dotfile_sync/libdotfilesync"

	"github.com/spf13/cobra"
)

var cmdSync *cobra.Command

func runSync(cmd *cobra.Command, args []string) error {
	verbose, err := rootCmd.Flags().GetBool("verbose")
	if err != nil {
		return err
	}

	if verbose == true {
	}
	fmt.Println("Running sync")
	libdotfilesync.Sync()
	return nil
}

func init() {
	cmdSync = &cobra.Command{
		Use:   "sync",
		Short: "Run a sync",
		RunE:  runSync,
	}

	rootCmd.AddCommand(cmdSync)
}
