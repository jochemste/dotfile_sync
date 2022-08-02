package dotfilesync

import (
	"errors"
	"fmt"
	"time"

	"github.com/jochemste/dotfile_sync/libdotfilesync"
)

func Sync(config *libdotfilesync.Config, configfile string) error {
	repo := libdotfilesync.NewRepo()
	var filemaps []*libdotfilesync.FileMap // variable to keep track of filemappings
	syncconfig := libdotfilesync.NewSyncConfig()
	config.Print()
	fmt.Println("Syncing...")

	// Clone the repository to memory
	err := repo.CloneRepo(config.UserSettings.Origin, config.UserSettings.AccessToken)
	if err != nil {
		return errors.New("Could not clone repo: " + err.Error())
	}

	// Check if Syncconfig file exists in repo
	if libdotfilesync.ExistsInFS(libdotfilesync.Syncfilepath) == false {
		// If it does not exist, create it and set the time back more than 10 years
		t := time.Now()
		t = t.AddDate(-10, -1, -1)
		syncconfig.SetLastSync(t)
		syncconfig.ToFS()
		fmt.Println("Syncconfig did not exist yet")
	} else {
		// If it does exit, load it
		fmt.Println("Syncconfig exists")
		syncconfig.FromFS()
	}
	syncconfig.Print()

	// Loop through the files to check for changes
	for _, f := range config.UserSettings.Files {

		// Create filemap for file
		fm := libdotfilesync.NewFileMap()
		err = fm.SetFilename(f)
		//fm.Origin = f
		//fm.FSPath, err = libdotfilesync.FindInFS(fm.GetOriginFilename())
		if err != nil && err.Error() == "Could not find file: File was not found" {
			fmt.Printf("%s does not exist in repo yet.\n", fm.GetOriginFilename())
		} else if err != nil {
			return errors.New("Could not find file " + fm.GetOriginFilename() + ": " + err.Error())
		}
		fm.FSTime = syncconfig.LastSync
		fm.Refresh()

		// Determine if the file has changed (either local or remote)
		isdiff, err := fm.HasChanged()
		if err != nil {
			return errors.New("Could not get diff: " + err.Error())
		}

		// If the file has changed
		if isdiff == true {
			// Update the file and determine if a commit is needed
			fm.Update(config.UserSettings.Mode)

			// If a commit is needed, commit the file to the repo
			if fm.NeedsCommit == true {
				err = repo.CommitToRepo(fm.GetFSFilename(), fm.Message)
				if err != nil {
					return errors.New("Failed to commit " + fm.GetOriginFilename() + ": " + err.Error())
				}
				syncconfig.SetLastSync(time.Now())
			}
		}
		filemaps = append(filemaps, fm)
	}

	// Commit syncconfig
	if libdotfilesync.SyncChanged == true {
		syncconfig.ToFS()
		err = repo.CommitToRepo(libdotfilesync.Syncfilepath, "Changed syncconfig")
		if err != nil {
			return errors.New("Failed to commit syncconfig: " + err.Error())
		}
	}

	config.DoNotChange.NrSync += 1
	config.DoNotChange.LastCheck = time.Now()
	err = config.ToFile()
	if err != nil {
		return errors.New("Could not write config file: " + err.Error())
	}

	// Push the repository to the remote
	err = repo.PushToRemote(config.UserSettings.AccessToken)
	if err != nil {
		return errors.New("Could not push to remote: " + err.Error())
	}

	return nil
}
