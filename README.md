#commandmocker

[![Build Status](https://secure.travis-ci.org/timeredbull/commandmocker.png?branch=master)](http://travis-ci.org/timeredbull/commandmocker)

commandmocker is a simple utility for tests in Go. It adds command with expected output to the path.

For example, if you want to mock the command "ssh", you can write a test that looks like this:

    import (
        "github.com/timeredbull/commandmocker"
        "testing"
    )

    func TestScreamsIfSSHFail(t *testing.T) {
        path, err := commandmocker.Add("ssh", "ssh: Could not resolve hostname myhost: nodename nor servname provided, or not known")
        if err != nil {
            t.Error(err)
            t.FailNow()
        }
        defer commandmocker.Remove(path)

        // write your test and expectations
    }
