// Copyright 2018 Shinichi MOTOKI. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
	"github.com/shimt/go-bufpool"
	"github.com/shimt/go-logif"
)

var (
	bufferPool = bufpool.NewBytesBufferPool(os.Getpagesize(), 0)
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

func (c *tmpleRuntime) makeRelPath(path string) (string, error) {
	p, err := filepath.Rel(c.base, path)
	if err != nil {
		return "", err
	}

	return p, nil
}

func (c *tmpleRuntime) makeTemplateName(path string) (string, error) {
	p, err := c.makeAbsPath(path)
	if err != nil {
		return "", errors.Wrap(err, "template name error")
	}

	p, err = c.makeRelPath(p)
	if err != nil {
		return "", errors.Wrap(err, "template name error")
	}

	return p, nil
}

func (c *tmpleRuntime) readFileToBuffer(buffer *bytes.Buffer, path string) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
	}()

	_, err = buffer.ReadFrom(f)
	return err
}

func (c *tmpleRuntime) processFile(path string, processor func(*bytes.Buffer) error) (err error) {
	b := bufferPool.Get()
	defer bufferPool.Put(b)

	if err = c.readFileToBuffer(b, path); err != nil {
		return err
	}

	return processor(b)
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
