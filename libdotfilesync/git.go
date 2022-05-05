package libdotfilesync

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jochemste/dotfile_sync/utils"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
)

type Gitter struct {
	Repository *git.Repository
	Storer     *memory.Storage
	Worktree   *git.Worktree
	//Filesystem FileSystem
}

func NewRepo() *Gitter {
	g := Gitter{}
	g.Storer = memory.NewStorage()
	//g.Filesystem = NewFS()
	return &g
}

func DetermMethod(origin string) (string, error) {
	if strings.HasPrefix(origin, "git@") {
		return "ssh", nil
	} else if strings.HasPrefix(origin, "http") {
		return "http", nil
	} else {
		return "", errors.New("Could not determine method")
	}
}

// Tests if a repository exists
func RepoExists(store_loc string) bool {
	res := utils.DirExists(store_loc)
	return res
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

func (g *Gitter) CloneRepo(origin string, store_loc string, secret ...string) error {
	NewFS()
	//fs := NewFS()
	/*method, err := DetermMethod(origin)
	if err != nil {
		return nil, err
	}*/

	//Delete the repo if it exists
	err := DeleteRepoIfExists(store_loc)
	if err != nil {
		return err
	}

	if len(secret) != 1 {
		// Git clone
		g.Repository, err = git.Clone(g.Storer, Filesystem.FS, &git.CloneOptions{
			URL:   origin,
			Depth: 5,
		})

		if err != nil {
			return err
		}
	} else {

		// Git clone
		g.Repository, err = git.Clone(g.Storer, Filesystem.FS, &git.CloneOptions{
			Auth: &http.BasicAuth{
				Username: "dotfile_sync", // This can be anything, except for an empty string
				Password: secret[0],
			},
			URL:   origin,
			Depth: 5,
		})

		if err != nil {
			return err
		}
	}

	g.Worktree, err = g.Repository.Worktree()

	return err
}

func (g *Gitter) CommitToRepo(filename string, message string) error {
	status, err := g.Worktree.Status()
	if err != nil {
		return err
	}
	fmt.Println(status)
	fmt.Printf("Committing %s\n", filename)

	//var worktree *git.Worktree
	//worktree, err := g.Repository.Worktree()
	//if err != nil {
	//	errors.New("Could not commit " + filename + ": " + err.Error())
	//}

	//worktree.Add(filename)
	g.Worktree.Add(filename)

	//worktree.Commit("AUTO: "+message, &git.CommitOptions{})
	commit, err := g.Worktree.Commit("AUTO: "+message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "dotfile_sync",
			Email: "none@none.test",
			When:  time.Now(),
		}})
	if err != nil {
		return err
	}

	status, err = g.Worktree.Status()
	if err != nil {
		return err
	}
	fmt.Println(status)

	obj, err := g.Repository.CommitObject(commit)
	fmt.Println(obj)

	return nil
}

func (g *Gitter) PushToRemote(secret string) error {
	fmt.Printf("Pushing to remote\n")
	err := g.Repository.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: "dotfile_sync", // This can be anything, except for an empty string
			Password: secret,
		},
	})

	if err != nil && err.Error() != "already up-to-date" {
		return errors.New("Push to remote failed: " + err.Error())
	} else if err != nil && err.Error() == "already up-to-date" {
		fmt.Println("Did not push to remote: already up-to-date")
	}

	return nil
}
