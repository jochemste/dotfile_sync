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

	if verbose == true {
	}
	fmt.Println("Running init")
	dotfilesync.Init("./test.toml")
	return nil
}

func init() {
	cmdInit = &cobra.Command{
		Use:   "init",
		Short: "Initialise",
		RunE:  runInit,
	}

	rootCmd.AddCommand(cmdInit)
}
