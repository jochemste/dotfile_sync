package libdotfilesync

import (
	//"os"
	"testing"

	"github.com/jochemste/dotfile_sync/utils"
)

func TestInitWorks(t *testing.T) {
	g := NewRepo()
	if g == nil {
		t.Errorf("NewRepo pointer is nil")
	}
}

func TestRepoClonePublic(t *testing.T) {
	g := NewRepo()
	err := g.CloneRepo("https://github.com/jochemste/dotfile_sync.git")
	if err != nil {
		t.Errorf("Public repo could not be cloned: %s", err)
	}
}

func TestRepoClonePublicNonExisting(t *testing.T) {
	g := NewRepo()
	err := g.CloneRepo("https://github.com/jochemste/dotfile_synchronothingness.git")
	if err == nil {
		t.Errorf("Non existing repo clone did not throw an error")
	}
}

func TestCommitToRepo(t *testing.T) {
	g := NewRepo()
	err := g.CloneRepo("https://github.com/jochemste/dotfile_sync.git")
	if err != nil {
		t.Errorf("Public repo could not be cloned: %s", err)
	}

	testfilename := "/tmp/test.test"
	err = utils.CreateFile(testfilename)
	if err != nil {
		t.Errorf("Could not create filename: %s", err)
	}

	err = g.CommitToRepo(testfilename, "testing")
	if err != nil {
		t.Errorf("Could not commit file: %s", err)
	}

	err = utils.RemoveFile(testfilename)
	if err != nil {
		t.Errorf("Could not remove filename: %s", err)
	}
}

func TestCommitToRepoNonExisting(t *testing.T) {
	g := NewRepo()
	err := g.CloneRepo("https://github.com/jochemste/dotfile_sync.git")
	if err != nil {
		t.Errorf("Public repo could not be cloned: %s", err)
	}

	testfilename := "/tmp/test.test"

	err = g.CommitToRepo(testfilename, "testing")
	if err != nil {
		t.Errorf("Could not commit file: %s", err)
	}
}
