package libdotfilesync

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
)

type Gitter struct {
	Repository *git.Repository
	Storer     *memory.Storage
	Worktree   *git.Worktree
}

func NewRepo() *Gitter {
	g := Gitter{}
	g.Storer = memory.NewStorage()
	return &g
}

// Clone a remote repository into memory
func (g *Gitter) CloneRepo(origin string, secret ...string) error {
	NewFS()

	var err error

	if len(secret) != 1 {
		// Git clone without a git token (public repos)
		g.Repository, err = git.Clone(g.Storer, Filesystem.FS, &git.CloneOptions{
			URL:   origin,
			Depth: 5,
		})

		if err != nil {
			return err
		}
	} else {

		// Git clone with a git token (private repos)
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
	fmt.Printf("Committing %s\n", filename)

	g.Worktree.Add(filename)

	_, err := g.Worktree.Commit("AUTO: "+message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "dotfile_sync",
			Email: "none@none.test",
			When:  time.Now(),
		}})
	if err != nil {
		return err
	}

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
