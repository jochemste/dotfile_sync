package libdotfilesync

import (
	//"os"
	"testing"
)

func TestRepoExists(t *testing.T) {
	existing := [2]string{"/tmp", "/etc"}
	nonexisting := [3]string{"/tmp/dotfile_syncc", "notanexistingrepo", "!@#%%43mksd"}

	for _, r := range existing {
		res := RepoExists(r)
		if res == false {
			t.Errorf("Repo %s does not appear to exist, but should.", r)
		}
	}

	for _, r := range nonexisting {
		res := RepoExists(r)
		if res == true {
			t.Errorf("Repo %s does appears to exist, but should not.", r)
		}
	}
}

/*
func TestGetRepoSSH(t *testing.T) {
	repo_url := "git@github.com:jochemste/dotfile_sync.git"
	privkey_file := os.Getenv("PWD") + "/../tmp/id_rsa"

	_, err := CloneRepo(repo_url, "/tmp/dotfilesync_test", privkey_file)
	if err != nil {
		t.Errorf("Git clone did not work: %s", err)
	}
}
*/

func TestCloneRepoHTTP(t *testing.T) {
	repo_url := "https://github.com/jochemste/dotfile_sync.git"

	_, err := CloneRepo(repo_url, "/tmp/dotfilesync_test")
	if err != nil {
		t.Errorf("Git clone did not work: %s", err)
	}

	exists := RepoExists("/tmp/dotfilesync_test")
	if exists == false {
		t.Errorf("Repository does not exist")
	}
}

func TestFindInFS(t *testing.T) {
	repo_url := "https://github.com/jochemste/dotfile_sync.git"
	_, err := CloneRepo(repo_url, "/tmp/dotfilesync_test")
	if err != nil {
		t.Errorf("Cloning failed: %s", err)
	}
	_, err = FindInFS(".gitignore")
	if err != nil {
		t.Errorf("Find in repository failed: %s", err)
	}

}
