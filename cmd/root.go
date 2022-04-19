package cmd

import (
	"log"

	"github.com/jochemste/dotfile_sync/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dotfilesync",
	Short: "Dotfile Sync is software to synchronise Linux configuration files (such as .bashrc, .emacs, etc) between devices",
	Long:  "",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().StringP("config", "c", utils.CONFDIR+"/config.toml", "Configuration file to be used")
}
