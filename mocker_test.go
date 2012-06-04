// Copyright 2012 commandmocker authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package commandmocker

import (
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
)

func TestAddFunctionReturnADirectoryThatIsInThePath(t *testing.T) {
	dir, err := Add("ssh", "success")
	defer Remove(dir)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	path := os.Getenv("PATH")
	if !strings.HasPrefix(path, dir) {
		t.Errorf("%s should be added to the first position in the path, but it was not.\nPATH: %s", dir, path)
	}
}

func TestAddFunctionShouldPutAnExecutableInTheReturnedDirectoryThatPrintsTheGivenOutput(t *testing.T) {
	dir, err := Add("ssh", "success")
	defer Remove(dir)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, err = os.Stat(path.Join(dir, "ssh"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	out, err := exec.Command("ssh").Output()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if string(out) != "success" {
		t.Errorf("should print success by running ssh, but printed %s", string(out))
	}
}

func TestRemoveFunctionShouldRemoveTheTempDirFromPath(t *testing.T) {
	dir, _ := Add("ssh", "success")
	err := Remove(dir)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	path := os.Getenv("PATH")
	if strings.HasPrefix(path, dir) {
		t.Errorf("%s should not be in the path, but it is.\nPATH: %s", dir, path)
	}
}

func TestRemoveFunctionShouldRemoveTheTempDirFromFileSystem(t *testing.T) {
	dir, _ := Add("ssh", "success")
	Remove(dir)
	_, err := os.Stat(dir)
	if err == nil || !os.IsNotExist(err) {
		t.Errorf("Directory %s should not exist, but it does.", dir)
	}
}

func TestShouldNotRemoveTheFirstItemWhenTheGivenDirectoryIsNotTheFirstInThePath(t *testing.T) {
	dir := os.TempDir() + "/blabla"
	err := Remove(dir)
	if err == nil || err.Error() != dir+" is not in $PATH" {
		t.Errorf("Should not be able to remove a directory that is not in $PATH")
	}
}

func TestShouldRemoveDirectoryFromArbitraryLocationInPath(t *testing.T) {
	dir, _ := Add("ssh", "success")
	path := os.Getenv("PATH")
	os.Setenv("PATH", "/:"+path)
	err := Remove(dir)
	path = os.Getenv("PATH")
	if err != nil || strings.Contains(path, dir) {
		t.Errorf("%s should not be in $PATH, but it is.", dir)
	}
}

func TestRemoveShouldReturnErrorIfTheGivenDirectoryDoesNotStartWithSlashTmp(t *testing.T) {
	err := Remove("/some/usr/bin")
	if err == nil || err.Error() != "Remove can only remove temporary directories, tryied to remove /some/usr/bin" {
		t.Error("Should not be able to remove non-temporary directories, but it was.")
	}
}
