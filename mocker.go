// Copyright 2012 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package commandmocker

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"syscall"
	"text/template"
)

var source = `#!/bin/sh -e

output="{{.}}"
echo -n "${output}"
`

func Add(name, output string) (tempdir string, err error) {
	tempdir = path.Join(os.TempDir(), "commandmocker")
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
	path := os.Getenv("PATH")
	path = tempdir + ":" + path
	err = os.Setenv("PATH", path)
	return
}

func Remove(tempdir string) error {
	path := os.Getenv("PATH")
	if strings.HasPrefix(path, tempdir) {
		path = path[len(tempdir)+1:]
		err := os.Setenv("PATH", path)
		if err != nil {
			return err
		}
		return os.RemoveAll(tempdir)
	}
	return errors.New(fmt.Sprintf("%s is not in $PATH", tempdir))
}
