// Copyright 2018 Shinichi MOTOKI. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// ---- kevArg ----

type kevArg string

func isKevArg(s string) bool {
	return strings.IndexByte(s, '=') > 0
}

func (a kevArg) parse() (key string, value string) {
	s := string(a)
	k := ""
	v := ""

	if i := strings.IndexByte(s, '='); i > 0 {
		k = s[0:i]
		v = s[i+1:]
	}

	return k, v
}

func (a kevArg) keyValue() (key string, value string) {
	return a.parse()
}

// ---- fileArg ----

type fileArg string

func isFileArg(s string) bool {
	return len(s) > 2 && s[0] == '@' && s[1] != '@'
}

func (a fileArg) parse() (fileType string, filePath string) {
	s := string(a)
	s = s[1:]

	if i := strings.IndexByte(s, ':'); i >= 0 {
		fileType = strings.ToUpper(s[0:i])
		s = s[i+1:]
	}

	filePath = s

	return fileType, filePath
}

func (a fileArg) readCloser() (readCloser io.ReadCloser, err error) {
	_, fp := a.parse()

	switch fp {
	case "STDIN":
		readCloser = ioutil.NopCloser(os.Stdin)
	default:
		readCloser, err = os.Open(fp)
		if err != nil {
			return nil, errors.Wrap(err, string(a))
		}
	}

	return readCloser, nil
}

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }

// NopWriteCloser returns a WriteCloser with a no-op Close method wrapping
// the provided Reader r.
func NopWriteCloser(w io.Writer) io.WriteCloser {
	return nopWriteCloser{w}
}

func (a fileArg) writeCloser() (writeCloser io.WriteCloser, err error) {
	_, fp := a.parse()

	switch fp {
	case "STDOUT":
		writeCloser = NopWriteCloser(os.Stdout)
	case "STDERR":
		writeCloser = NopWriteCloser(os.Stderr)
	default:
		writeCloser, err = os.OpenFile(fp, os.O_CREATE, os.FileMode(0644))
		if err != nil {
			return nil, errors.Wrap(err, string(a))
		}
	}

	return writeCloser, nil
}

func (a fileArg) blob() (blob []byte, err error) {
	rc, err := a.readCloser()
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rc.Close()
	}()

	blob, err = ioutil.ReadAll(rc)
	return blob, errors.Wrap(err, string(a))
}

// --- stringArg ----

type stringArg string

func (a stringArg) String() string {
	s := string(a)
	l := len(s)

	if l == 0 {
		return s
	}

	if strings.Contains("@!", s[0:1]) {
		if l >= 2 && s[0] == s[1] {
			s = s[1:]
		}
	}

	return s
}
