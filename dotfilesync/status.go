package dotfilesync

import (
	"errors"
	"fmt"
	//"time"

	"github.com/jochemste/dotfile_sync/libdotfilesync"
)

func Status(config *libdotfilesync.Config, configfile string) error {

	// Get repo
	repo := libdotfilesync.NewRepo()
	//var filemaps []*libdotfilesync.FileMap // variable to keep track of filemappings
	syncconfig := libdotfilesync.NewSyncConfig()
	config.Print()
	fmt.Println("Getting status...")
	// Clone the repo to memory
	err := repo.CloneRepo(config.UserSettings.Origin, config.UserSettings.AccessToken)
	if err != nil {
		return errors.New("Could not clone repo: " + err.Error())
	}

	// Check if syncconfig file exists in repo
	if libdotfilesync.ExistsInFS(libdotfilesync.Syncfilepath) == false {
		return errors.New("Could not get status, since syncfile does not exists in remote repo")
	} else {
		// If it does exit, load it
		fmt.Println("Syncconfig exists")
		syncconfig.FromFS()
	}
	syncconfig.Print()

	// Compare repo file datetimes with local file datetimes
	// Loop through the files
	for _, f := range config.UserSettings.Files {

		// Create filemap for file
		fm := libdotfilesync.NewFileMap()
		err = fm.SetFilename(f)
		fm.FSTime = syncconfig.LastSync

		if err != nil && err.Error() == "Could not find file: File was not found" {
			fmt.Printf("%s does not exist in repo yet.\n", fm.GetOriginFilename())
		} else if err != nil {
			return errors.New("Could not find file " + fm.GetOriginFilename() + ": " + err.Error())
		}

		// Determine if the file has changed (either local or remote)
		isdiff, err := fm.HasChanged()
		if err != nil {
			return errors.New("Could not get diff: " + err.Error())
		}

		if isdiff == true {
			fmt.Printf("\t- %s has changed\n", fm.GetOriginFilename())
			if fm.LocalIsMoreRecent() {
				fmt.Printf("\t\tLocal is more Recent\n")
				fmt.Printf("\t\t%s - %s\n", fm.OriginTime, fm.FSTime)
			} else if fm.FMIsMoreRecent() {
				fmt.Printf("\t\tRemote is more Recent\n")
				fmt.Printf("\t\tOG:%s - REPO:%s\n", fm.OriginTime, fm.FSTime)
			} else {
				fmt.Printf("\t\tDatetimes are equal\n")
				fmt.Printf("\t\tOG:%s - REPO:%s\n", fm.OriginTime, fm.FSTime)
			}

		} else {
			fmt.Printf("\t- %s has not changed\n", fm.GetOriginFilename())
		}
	}

	return nil
}
