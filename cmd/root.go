package cmd

import (
	//"fmt"
	"log"
	"os"

	//"github.com/jochemste/dotfile_sync/libdotfilesync"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dotfilesync",
	Short: "Dotfile Sync is software to synchronise Linux configuration files (such as .bashrc, .emacs, etc) between devices",
	Long:  "",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(os.Stderr, err)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
}
