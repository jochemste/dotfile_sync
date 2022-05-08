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

	err := repo.CloneRepo(config.UserSettings.Origin, config.DoNotChange.Store, config.UserSettings.AccessToken)
	if err != nil {
		return errors.New("Could not clone repo: " + err.Error())
	}

	// Check if Syncconfig file exists in repo
	if libdotfilesync.ExistsInFS(libdotfilesync.Syncfilepath) == false {
		t := time.Now()
		t = t.AddDate(-10, -1, -1)
		fmt.Println(t)
		fmt.Println(time.Now())
		syncconfig.SetLastSync(t)
		syncconfig.ToFS()
		fmt.Println("Syncconfig did not exist yet")
	} else {
		fmt.Println("Syncconfig did not exist yet")
		syncconfig.FromFS()
	}
	syncconfig.Print()

	// Loop through the files to check for changes
	for _, f := range config.UserSettings.Files {
		fm := libdotfilesync.NewFileMap()
		fm.Origin = f
		fm.FSPath, err = libdotfilesync.FindInFS(fm.GetOriginFilename())
		if err != nil {
			return errors.New("Could not find file " + fm.GetOriginFilename() + ": " + err.Error())
		}
		fm.FSTime = syncconfig.LastSync
		fm.Refresh()

		isdiff, err := fm.HasChanged()
		if err != nil {
			return errors.New("Could not get diff: " + err.Error())
		}

		fmt.Printf("%s has changed: %t\n", fm.GetOriginFilename(), isdiff)
		if isdiff == true {
			fm.Update(config.UserSettings.Mode)
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

	err = repo.PushToRemote(config.UserSettings.AccessToken)
	if err != nil {
		return errors.New("Could not push to remote: " + err.Error())
	}

	return nil
}
