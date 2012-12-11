// Copyright 2012 commandmocker authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package commandmocker

import (
	"bytes"
	"io/ioutil"
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

func TestOutputFunctionReturnsOutputOfExecutedCommand(t *testing.T) {
	dir, err := Add("ssh", "$*")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer Remove(dir)
	_, err = exec.Command("ssh", "foo", "bar").Output()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	out := Output(dir)
	if out != "foo bar" {
		t.Errorf("Output function should return output of ssh command execution, got '%s'", out)
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
	pathMutex.Lock()
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

func TestRanCheckIfTheDotRanFileExists(t *testing.T) {
	dir, err := Add("ls", "bla")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer Remove(dir)
	p := path.Join(dir, ".ran")
	f, err := os.Create(p)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer f.Close()
	table := map[string]bool{
		dir:            true,
		"/tmp/bla/bla": false,
		"/home/blabla": false,
	}
	for input, expected := range table {
		got := Ran(input)
		if got != expected {
			t.Errorf("Ran on %s?\nExpected: %q.\nGot: %q.", input, expected, got)
		}
	}
}

func TestErrorGeneratesTheFileThatReturnsExitStatusCode(t *testing.T) {
	var (
		content, p string
		b          []byte
	)
	dir, err := Error("ssh", "bla", 1)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer Remove(dir)
	p = path.Join(dir, "ssh")
	b, err = ioutil.ReadFile(p)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	content = string(b)
	if !strings.Contains(content, "exit 1") {
		t.Errorf(`Did not find "exit 1" in the generated file. Content: %s`, content)
	}
}

func TestErrorReturnsOutputInStderr(t *testing.T) {
	dir, err := Error("ssh", "ble", 42)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer Remove(dir)
	cmd := exec.Command("ssh")
	var b bytes.Buffer
	cmd.Stderr = &b
	err = cmd.Run()
	if err == nil {
		t.Error(err)
		t.FailNow()
	}
	if string(b.String()) != "ble" {
		t.Errorf("should print ble running ssh, but printed %s", b.String())
	}
}

func TestMultipleCallsAppendToOutput(t *testing.T) {
	dir, err := Add("ssh", "ble")
	if err != nil {
		t.Fatal(err)
	}
	defer Remove(dir)
	err = exec.Command("ssh").Run()
	if err != nil {
		t.Fatal(err)
	}
	err = exec.Command("ssh").Run()
	if err != nil {
		t.Fatal(err)
	}
	got := Output(dir)
	if got != "bleble" {
		t.Errorf("Output(%q): Want %q. Got %q.", dir, "bleble", got)
	}
}
