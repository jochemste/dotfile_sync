package cmd

import (
	//"log"
	"fmt"

	"github.com/jochemste/dotfile_sync/dotfilesync"

	"github.com/spf13/cobra"
)

var cmdInit *cobra.Command

func runInit(cmd *cobra.Command, args []string) error {
	verbose, err := rootCmd.Flags().GetBool("verbose")
	if err != nil {
		return err
	}
	configfile, err := rootCmd.Flags().GetString("config")
	if err != nil {
		return err
	}

	interactive, err := cmdInit.Flags().GetBool("interactive")
	if err != nil {
		return err
	}

	if verbose == true {
	}
	fmt.Println("Running init")
	err = dotfilesync.Init(interactive, configfile)
	return err
}

func init() {
	cmdInit = &cobra.Command{
		Use:   "init",
		Short: "Initialise",
		RunE:  runInit,
	}

	cmdInit.Flags().BoolP("interactive", "i", false, "Initialise interactively")

	rootCmd.AddCommand(cmdInit)
}
