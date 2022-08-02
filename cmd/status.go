package cmd

import (
	"fmt"

	"github.com/jochemste/dotfile_sync/dotfilesync"
	"github.com/jochemste/dotfile_sync/libdotfilesync"

	"github.com/spf13/cobra"
)

var cmdStatus *cobra.Command

/* Get the status of the repo and the files in it */
func runStatus(cmd *cobra.Command, args []string) error {
	verbose, err := rootCmd.Flags().GetBool("verbose")
	if err != nil {
		return err
	}

	if verbose == true {
	}
	fmt.Println("Running status")

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

	err = dotfilesync.Status(config, configfile)
	return err
}

func init() {
	cmdStatus = &cobra.Command{
		Use:   "status",
		Short: "Run a status check",
		RunE:  runStatus,
	}

	rootCmd.AddCommand(cmdStatus)
}
