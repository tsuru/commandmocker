// Copyright 2012 commandmocker authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// commandmocker is a simple utility for tests in Go. It adds command with
// expected output to the path.
//
// For example, if you want to mock the command "ssh", you can write a test that
// looks like this:
//
//     import (
//         "github.com/timeredbull/commandmocker"
//         "testing"
//     )
//
//     func TestScreamsIfSSHFail(t *testing.T) {
// 		msg := "ssh: Could not resolve hostname myhost: nodename nor servname provided, or not known"
//         path, err := commandmocker.Add("ssh", msg)
//         if err != nil {
//             t.Error(err)
//             t.FailNow()
//         }
//         defer commandmocker.Remove(path)
//
//         // write your test and expectations
//     }
//
// Please notice that commandmocker is intended for testing only.
package commandmocker
