// Copyright 2018 Shinichi MOTOKI. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
	"github.com/shimt/go-simplecli"
)

var cli = simplecli.NewCLI()
var cliInstruction = struct {
	in  string
	out string
}{}

func init() {
	err := cli.Initialize()
	cli.Exit1IfError(err)

	cli.CommandLine.StringVarP(&cliInstruction.in, "in", "i", "@STDIN", "template")
	cli.CommandLine.StringVarP(&cliInstruction.out, "out", "o", "@STDOUT", "output")
}

func parseTemplate(tc *tmpleContext) (string, *template.Template) {
	var (
		err error
		in  io.ReadCloser
		fp  string
		fd  string
	)

	fn := cliInstruction.in

	if isFileArg(fn) {
		f := fileArg(fn)
		in, err = f.readCloser()
		cli.Exit1IfError(err)

		fd, err = os.Getwd()
		cli.Exit1IfError(err)
	} else {
		fp = fn
		in, err = os.Open(fp)
		cli.Exit1IfError(err)

		fd, err = filepath.Abs(filepath.Dir(fp))
		cli.Exit1IfError(err)
	}
	defer in.Close()

	b, err := ioutil.ReadAll(in)
	cli.Exit1IfError(err)

	t, err := tc.newTemplate(fn).Parse(string(b))
	cli.Exit1IfError(err)

	return fd, t
}

func mergeData(dst map[string]interface{}, src map[string]interface{}) {
	for n, v := range src {
		dst[n] = v
	}
}

func processFile(s string) (data map[string]interface{}, args []interface{}) {
	_, _, fo, err := getFileObject(s, "TEXT")
	cli.Exit1IfError(errors.Wrap(err, s))

	switch v := fo.(type) {
	case map[string]interface{}:
		data = v
	case []interface{}:
		args = v
	default:
		args = []interface{}{v}
	}

	return data, args
}

func openOutput() (out io.WriteCloser) {
	var (
		err error
	)

	fn := cliInstruction.out

	if isFileArg(fn) {
		out, err = fileArg(fn).writeCloser()
		cli.Exit1IfError(err)
	} else {
		out, err = os.OpenFile(fn, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, os.FileMode(0644))
		cli.Exit1IfError(err)
	}

	return out
}

func main() {
	err := cli.Setup()
	cli.Exit1IfError(err)

	cli.StartProfile()
	defer cli.StopProfile()

	data := map[string]interface{}{}
	args := make([]interface{}, 0, cli.CommandLine.NArg())

	data["Environ"] = environMap()

	for _, v := range cli.CommandLine.Args() {
		switch {
		case isKevArg(v):
			nevN, nevV := kevArg(v).keyValue()
			switch {
			case isFileArg(nevV):
				_, _, o, e := getFileObject(nevV, "TEXT")
				cli.Exit1IfError(e)
				data[nevN] = o
			default:
				data[nevN] = stringArg(nevV).String()
			}
		case isFileArg(v):
			d, a := processFile(v)
			mergeData(data, d)
			args = append(args, a...)
		default:
			args = append(args, stringArg(v).String())
		}
	}

	data["Arguments"] = args

	tc := &tmpleContext{log: cli.Log}

	fd, tmpl := parseTemplate(tc)

	b := &bytes.Buffer{}

	err = tc.Execute(b, fd, tmpl, data)
	cli.Exit1IfError(err)

	out := openOutput()
	defer out.Close()

	_, err = b.WriteTo(out)
	cli.Exit1IfError(err)

	cli.Exit(0)
}
