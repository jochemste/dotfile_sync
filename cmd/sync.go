package cmd

import (
	"fmt"

	"github.com/jochemste/dotfile_sync/dotfilesync"
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

	// Getting config (file should be cli parameter with default value)
	config := libdotfilesync.NewConfig()
	configfile, err := rootCmd.Flags().GetString("config")
	if err != nil {
		return err
	}
	err = config.FromFile(configfile)
	if err != nil {
		return err
	}

	err = dotfilesync.Sync(config, configfile)
	return err
}

func init() {
	cmdSync = &cobra.Command{
		Use:   "sync",
		Short: "Run a sync",
		RunE:  runSync,
	}

	rootCmd.AddCommand(cmdSync)
}
