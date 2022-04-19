package dotfilesync

import (
	"fmt"
	//"log"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/jochemste/dotfile_sync/libdotfilesync"
	"github.com/jochemste/dotfile_sync/utils"
	//"github.com/pelletier/go-toml/v2"
)

func initNonInteractive(file_path string) error {
	fmt.Println("Initialising")

	homedir := os.Getenv("HOME") + "/"
	files_init := []string{homedir + ".bashrc", homedir + ".emacs", homedir + ".profile"}
	var files []string

	for _, f := range files_init {
		if utils.FileExists(f) {
			files = append(files, f)
		}
	}

	//Initialise configuration and write to file
	config := libdotfilesync.NewConfig()
	config.UserSettings.Origin = "http://example.com/replace_with_repository_url"
	config.UserSettings.Files = files
	config.UserSettings.AccessToken = "REPLACE_WITH_ACCESSTOKEN"
	config.DoNotChange.LastCheck = time.Now()
	config.DoNotChange.NrSync = 0
	config.DoNotChange.Store = "/tmp/dotfile_sync"

	if err := config.ToFile(file_path); err != nil {
		return err
	}

	fmt.Println("New configuration file created here: "+file_path, "\nMake sure to enter the correct paths to the files you would like to synchronise in there")
	return nil
}

func initInteractive(file_path string) error {
	var origin string
	var accesstoken string

	if utils.FileExists(file_path) {
		var choice string
		fmt.Println("File exists. Do you want to overwrite it? [y/n]")
		fmt.Scanln(&choice)
		if choice == "y" || choice == "Y" {
			//continue
		} else {
			fmt.Println("Okay, exiting")
			os.Exit(1)
		}
	}

	fmt.Println("Initialising")

	fmt.Println("Enter repository origin")
	fmt.Scanln(&origin)

	if strings.HasPrefix(origin, "git@") {
		return errors.New("SSH is not supported at the moment")
	} else if !strings.HasPrefix(origin, "http") {
		return errors.New("Could not determine URL, please make sure the format is a valid URL")
	}

	fmt.Println("Enter Personal Access Token")
	fmt.Scanln(&accesstoken)

	homedir := os.Getenv("HOME") + "/"
	files_init := []string{homedir + ".bashrc", homedir + ".emacs", homedir + ".profile"}
	var files []string

	for _, f := range files_init {
		if utils.FileExists(f) {
			files = append(files, f)
		}
	}

	//Initialise configuration and write to file
	config := libdotfilesync.NewConfig()
	config.UserSettings.Origin = "http://example.com/replace_with_repository_url"
	config.UserSettings.Files = files
	config.UserSettings.AccessToken = "REPLACE_WITH_ACCESSTOKEN"
	config.DoNotChange.LastCheck = time.Now()
	config.DoNotChange.NrSync = 0

	if err := config.ToFile(file_path); err != nil {
		return err
	}

	fmt.Println("New configuration file created here: "+file_path,
		"\nMake sure to enter the correct paths to the files you would like to synchronise in there")
	return nil
}

// Initialise Dotfile Sync. Creates a configuration file in the provided path, or uses the
// default path of none are provided. Has the option to interactively prompt the user
// for information or to non-interactively create the config.toml file and allow the user
// to modify this later.
func Init(interactive bool, file_path ...string) error {
	var fp string
	if len(file_path) == 0 {
		fp = utils.CONFDIR + "/config.toml"
	} else {
		fp = file_path[0]
	}

	if utils.DirExists(utils.CONFDIR) == false {
		err := utils.Mkdir(utils.CONFDIR)
		if err != nil {
			return err
		}
		stats, _ := os.Stat(utils.CONFDIR)
		fmt.Println(stats.Mode())
		err = utils.SetDefPermission(utils.CONFDIR)
		if err != nil {
			return err
		}
		stats, _ = os.Stat(utils.CONFDIR)
		fmt.Println(stats.Mode())
	}

	if interactive == true {
		initInteractive(fp)
	} else {
		initNonInteractive(fp)
	}

	return nil
}
