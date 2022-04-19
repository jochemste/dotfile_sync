package dotfilesync

import (
	"errors"
	"fmt"

	"github.com/jochemste/dotfile_sync/libdotfilesync"
)

func Sync(config *libdotfilesync.Config) error {
	fmt.Println("Syncing...")
	config.Print()

	_, err := libdotfilesync.CloneRepo(config.UserSettings.Origin, config.DoNotChange.Store, config.UserSettings.AccessToken)
	if err != nil {
		return errors.New("Could not clone repo: " + err.Error())
	}

	file, err := libdotfilesync.FindInFS(".emacs")
	if err != nil {
		return errors.New("Could not find file " + file + ": " + err.Error())
	}

	fmt.Println(file)

	return nil
}
