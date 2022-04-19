package libdotfilesync

import (
	"errors"
	"fmt"
	"os"
	"strings"
	//"time"

	"github.com/jochemste/dotfile_sync/utils"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	//"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
)

var filesystem = memfs.New()

func DetermMethod(origin string) (string, error) {
	if strings.HasPrefix(origin, "git@") {
		return "ssh", nil
	} else if strings.HasPrefix(origin, "http") {
		return "http", nil
	} else {
		return "", errors.New("Could not determine method")
	}
}

// Tests if a directory exists (although it will also work for files)
func RepoExists(store_loc string) bool {
	res := utils.DirExists(store_loc)

	if res == true {
		return true
	} else {
		return false
	}
}

// Delete the repository if it exists
func DeleteRepoIfExists(store_loc string) error {
	exists := RepoExists(store_loc)

	if exists == true {
		err := os.RemoveAll(store_loc)
		return err
	}
	return nil
}

func CloneRepo(origin string, store_loc string, secret ...string) (*git.Repository, error) {
	method, err := DetermMethod(origin)
	if err != nil {
		return nil, err
	}

	//Delete the repo if it exists
	err = DeleteRepoIfExists(store_loc)
	if err != nil {
		return nil, err
	}

	// HTTP method
	if method == "http" {
		if len(secret) != 1 {
			// Git clone
			repo, err := git.Clone(memory.NewStorage(), filesystem, &git.CloneOptions{
				URL:   origin,
				Depth: 5,
			})
			/*
				repo, err := git.PlainClone(store_loc, false, &git.CloneOptions{
					URL:   origin,
					Depth: 5,
				})*/
			return repo, err
		} else {

			// Git clone
			repo, err := git.Clone(memory.NewStorage(), filesystem, &git.CloneOptions{
				Auth: &http.BasicAuth{
					Username: "dotfile_sync", // This can be anything, except for an empty string
					Password: secret[0],
				},
				URL:   origin,
				Depth: 5,
			})
			/*
				repo, err := git.PlainClone(store_loc, false, &git.CloneOptions{
					Auth: &http.BasicAuth{
						Username: "dotfile_sync", // This can be anything, except for an empty string
						Password: secret[0],
					},
					URL:   origin,
					Depth: 5,
				})*/
			return repo, err
		}

		// SSH method
	} else if method == "ssh" {

		// Make sure a private key file is available
		if len(secret) != 1 {
			return nil, errors.New("Private key file is required when using SSH")
		} else if utils.FileExists(secret[0]) == false {
			return nil, errors.New("Private key file " + secret[0] + " does not exist")
		}
		privkey_file := secret[0]

		// Get public key from private key file
		publicKeys, err := ssh.NewPublicKeysFromFile("git", privkey_file, "")
		if err != nil {
			return nil, err
		}

		// Git clone
		repo, err := git.Clone(memory.NewStorage(), filesystem, &git.CloneOptions{
			Auth: publicKeys,
			URL:  origin,
		})
		/*
			repo, err := git.PlainClone(store_loc, false, &git.CloneOptions{
				Auth: publicKeys,
				URL:  origin,
			})*/

		return repo, err

		// UNKNOWN method
	} else {
		return nil, errors.New("Method not recognised, only http or ssh are allowed")
	}
}

// Find in filesystem
func FindInFS(filename string) (string, error) {
	files, err := filesystem.ReadDir("/")
	fmt.Printf("%T\n", files)
	if err != nil {
		return "", errors.New("Could not find file: " + err.Error())
	}

	res, err := findInFSDir(filename, files, "")
	if err != nil {
		return "", errors.New("Could not find file: " + err.Error())
	}

	return res, nil
}

// Recursive find in directory function (private)
func findInFSDir(filename string, files []os.FileInfo, dirname string) (string, error) {

	for _, file := range files {
		fmt.Printf("Filename: %s\n", dirname+"/"+file.Name())
		if file.IsDir() == true {
			f, err := filesystem.ReadDir(file.Name())
			if err != nil {
				return "", err
			}
			res, err := findInFSDir(filename, f, dirname+"/"+file.Name())
			if res != "" {
				return res, err
			}
		}

		if file.Name() == filename {
			return dirname + "/" + file.Name(), nil
		}
	}

	return "", errors.New("File was not found")
}
