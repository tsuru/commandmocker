// Copyright 2012 commandmocker authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package commandmocker

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"syscall"
	"text/template"
	"time"
)

var source = `#!/bin/sh -e

output="{{.}}"
echo -n "${output}"
`

var pathMutex sync.Mutex

// Add creates a temporary directory containing an executable file named "name"
// that prints "output" when executed. It also adds the temporary directory to
// the first position of $PATH.
//
// It returns the temporary directory path (for future removing, using the
// Remove function) and an error if any happen.
func Add(name, output string) (tempdir string, err error) {
	tempdir = path.Join(os.TempDir(), "commandmocker+" + time.Now().Format("20060102150405999999999"))
	_, err = os.Stat(tempdir)
	for !os.IsNotExist(err) {
		tempdir = path.Join(os.TempDir(), "commandmocker+" + time.Now().Format("20060102150405999999999"))
		_, err = os.Stat(tempdir)
	}
	err = os.MkdirAll(tempdir, 0777)
	if err != nil {
		return
	}
	f, err := os.OpenFile(path.Join(tempdir, name), syscall.O_WRONLY|syscall.O_CREAT|syscall.O_TRUNC, 0755)
	if err != nil {
		return
	}
	defer f.Close()
	t, err := template.New(name).Parse(source)
	if err != nil {
		return
	}
	err = t.Execute(f, output)
	if err != nil {
		return
	}
	pathMutex.Lock()
	defer pathMutex.Unlock()
	path := os.Getenv("PATH")
	path = tempdir + ":" + path
	err = os.Setenv("PATH", path)
	return
}

// Remove removes the tempdir from $PATH and from file system.
//
// This function is intended only to undo what Add does. It returns error if
// the given tempdir is not a temporary directory.
func Remove(tempdir string) error {
	if !strings.HasPrefix(tempdir, os.TempDir()) {
		return errors.New("Remove can only remove temporary directories, tryied to remove " + tempdir)
	}
	pathMutex.Lock()
	path := os.Getenv("PATH")
	index := strings.Index(path, tempdir)
	if index < 0 {
		pathMutex.Unlock()
		return errors.New(fmt.Sprintf("%s is not in $PATH", tempdir))
	}
	path = path[:index] + path[index+len(tempdir)+1:]
	err := os.Setenv("PATH", path)
	pathMutex.Unlock()
	if err != nil {
		return err
	}
	return os.RemoveAll(tempdir)
}
