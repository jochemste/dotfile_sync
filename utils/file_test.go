package utils

import (
	"os"
	"testing"
)

func TestDirFileSplit(t *testing.T) {
	files := []string{"/etc/passwd", "./someotherfile", ".gitignore", "../../../.emacs"}
	expected_dirs := []string{"/etc/", "./", "", "../../../"}
	expected_files := []string{"passwd", "someotherfile", ".gitignore", ".emacs"}

	for i, f := range files {
		dir, file := DirFileSplit(f)
		if expected_dirs[i] != dir {
			t.Errorf("%s is not the same as the expected dir %s", dir, expected_dirs[i])
		}
		if expected_files[i] != file {
			t.Errorf("%s is not the same as the expected file %s", file, expected_files[i])
		}
	}
}

func TestPathExists(t *testing.T) {
	existingdirs := [3]string{"/etc", "/etc/",
		os.Getenv("HOME")}
	existingfiles := [3]string{"/etc/passwd",
		os.Getenv("HOME") + "/.bashrc",
		os.Getenv("HOME") + "/.profile"}
	nonexisting := [3]string{"/somerandomfilenamethatshouldnotexist",
		"idontevenknow what to put here",
		"123453rkmnfrls"}

	for _, f := range existingdirs {
		res, pt, err := PathExists(f)
		if err != nil {
			t.Errorf("Path %s was not found while it should exist: %s, %d", f, err, pt)
		}
		if res == false {
			t.Errorf("Path %s was not found, while it should exist", f)
		}
		if pt != DIR {
			t.Errorf("Path %s was not found to be a dir, while it is", f)
		}
	}

	for _, f := range existingfiles {
		res, pt, err := PathExists(f)
		if err != nil {
			t.Errorf("Path %s was not found while it should exist: %s, %d", f, err, pt)
		}
		if res == false {
			t.Errorf("Path %s was not found, while it should exist", f)
		}
		if pt != FILE {
			t.Errorf("Path %s was not found to be a file, while it is", f)
		}
	}

	for _, f := range nonexisting {
		res, _, _ := PathExists(f)
		if res == true {
			t.Errorf("Path %s was found, while it does not exist", f)
		}
	}
}

func TestFileExists(t *testing.T) {
	existing := [3]string{"/etc/passwd",
		os.Getenv("HOME") + "/.bashrc",
		os.Getenv("HOME") + "/.profile"}
	nonexisting := [3]string{"/somerandomfilenamethatshouldnotexist",
		"idontevenknow what to put here",
		"123453rkmnfrls"}

	for _, f := range existing {
		res := FileExists(f)
		if res == false {
			t.Errorf("File %s was not found, while it should exist", f)
		}
	}

	for _, f := range nonexisting {
		res := FileExists(f)
		if res == true {
			t.Errorf("File %s was found, while it does not exist", f)
		}
	}
}

func TestDirExists(t *testing.T) {
	existing := [3]string{"/etc", "/etc/",
		os.Getenv("HOME")}
	existingfiles := [3]string{"/etc/passwd",
		os.Getenv("HOME") + "/.bashrc",
		os.Getenv("HOME") + "/.profile"}
	nonexisting := [3]string{"/somerandomfilenamethatshouldnotexist",
		"idontevenknow what to put here",
		"123453rkmnfrls"}

	for _, f := range existing {
		res := DirExists(f)
		if res == false {
			t.Errorf("Dir %s was not found, while it should exist", f)
		}
	}

	for _, f := range existingfiles {
		res := DirExists(f)
		if res == true {
			t.Errorf("Dir %s was found, while it is a file", f)
		}
	}

	for _, f := range nonexisting {
		res := DirExists(f)
		if res == true {
			t.Errorf("Dir %s was found, while it does not exist", f)
		}
	}
}

func TestGetLastModified(t *testing.T) {
	_, err := GetLastModified("/etc/passwd")
	if err != nil {
		t.Errorf("Error while getting modtime")
	}
}

func TestIsMoreRecent(t *testing.T) {
	f1 := "/etc/passwd"
	f2 := os.Getenv("HOME") + "/.bashrc"
	_, err := IsMoreRecent(f1, f2)
	if err != nil {
		t.Errorf("Error while comparing modtimes from %s and %s", f1, f2)
	}
}

func TestIsDiff(t *testing.T) {
	f1 := "/etc/passwd"
	f2 := os.Getenv("HOME") + "/.bashrc"
	diffs, err := IsDiff(f1, f2)
	if err != nil {
		t.Errorf("Error while comparing files")
	}
	if diffs == false {
		t.Errorf("Files are not different, while they should be")
	}

	diffs, err = IsDiff(f1, f1)
	if err != nil {
		t.Errorf("Error while comparing files")
	}
	if diffs == true {
		t.Errorf("Files are different, while they should not be")
	}
}
