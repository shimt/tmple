// Copyright 2018 Shinichi MOTOKI. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
	"github.com/shimt/go-logif"
)

type tmpleRuntime struct {
	tmpl *template.Template

	base string
	dir  dirStack

	data     map[string]interface{}
	fullpath map[string]string

	log logif.LeveledLogger
}

func (c *tmpleRuntime) newTemplate(name string) *template.Template {
	return template.New(name).Funcs(c.funcMap())
}

func (c *tmpleRuntime) makeAbsPath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	p := filepath.Clean(filepath.Join(c.dir.getCwd(), path))

	return p, nil
}

func (c *tmpleRuntime) makeTemplateName(path string) (string, error) {
	p, err := c.makeAbsPath(path)
	if err != nil {
		return "", errors.Wrap(err, "template name error")
	}

	p, err = filepath.Rel(c.base, p)
	if err != nil {
		return "", errors.Wrap(err, "template name error")
	}

	return p, nil
}

func (c *tmpleRuntime) Execute(out io.Writer, base string, tmpl *template.Template, data map[string]interface{}) (err error) {
	c.tmpl = tmpl
	c.data = data
	c.base = base
	c.fullpath = map[string]string{}

	return c.dir.run(base, func() error {
		return tmpl.Execute(out, data)
	})
}
