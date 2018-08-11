// Copyright 2018 Shinichi MOTOKI. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"path/filepath"
	"sort"
	"text/template"

	"github.com/pkg/errors"
)

var (
	errUnsupportedArgumentType = errors.New("unsupported argument type")
)

func (c *tmpleRuntime) funcMap() template.FuncMap {
	return template.FuncMap{
		"glob":            c.tfGlob,
		"includeFile":     c.tfIncludeFile,
		"includeTemplate": c.tfIncludeTemplate,
		"includeTextFile": c.tfIncludeTextFile,
	}
}

func (c *tmpleRuntime) tfGlob(glob string) ([]string, error) {
	m, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}

	sort.Strings(m)

	return m, nil
}

func argToStringSlice(args []interface{}) ([]string, error) {
	l := make([]string, 0, len(args))

	for _, a := range args {
		switch v := a.(type) {
		case string:
			l = append(l, v)
		case []string:
			l = append(l, v...)
		default:
			return nil, errUnsupportedArgumentType
		}
	}

	return l, nil
}

func (c *tmpleRuntime) tfIncludeFile(args ...interface{}) (string, error) {
	paths, err := argToStringSlice(args)
	if err != nil {
		return "", err
	}

	r := bufferPool.Get().(*bytes.Buffer)
	r.Reset()
	defer bufferPool.Put(r)

	for _, p := range paths {
		err = c.readFileToBuffer(r, p)
		if err != nil {
			return "", err
		}
	}

	return r.String(), nil
}

func (c *tmpleRuntime) tfIncludeTextFile(args ...interface{}) (string, error) {
	paths, err := argToStringSlice(args)
	if err != nil {
		return "", err
	}

	r := bufferPool.Get().(*bytes.Buffer)
	r.Reset()
	defer bufferPool.Put(r)

	for _, p := range paths {
		err = c.readFileToBuffer(r, p)
		if err != nil {
			return "", err
		}

		b := r.Bytes()
		if c := b[len(b)-1]; c != '\r' && c != '\n' {
			r.WriteByte('\n')
		}
	}

	return r.String(), nil
}

func (c *tmpleRuntime) tfIncludeTemplate(args ...interface{}) (string, error) {
	paths, err := argToStringSlice(args)
	if err != nil {
		return "", err
	}

	r := bufferPool.Get().(*bytes.Buffer)
	r.Reset()
	defer bufferPool.Put(r)

	for _, p := range paths {
		fp, err := c.makeAbsPath(p)
		if err != nil {
			return "", err
		}

		tn, err := c.makeTemplateName(fp)
		if err != nil {
			return "", err
		}

		err = c.processFile(fp, func(tb *bytes.Buffer) error {
			c.log.Debugf("tfIncludeTemplate: load template %s (%s)", tn, fp)

			_, err = c.tmpl.New(tn).Parse(tb.String())

			return err
		})

		c.fullpath[tn] = fp

		err = c.dir.run(filepath.Dir(fp), func() error {
			c.log.Debugf("tfIncludeTemplate: change working directory %s", c.dir.getCwd())
			c.log.Debugf("tfIncludeTemplate: execute template %s", tn)

			return c.tmpl.ExecuteTemplate(r, p, c.data)
		})
		if err != nil {
			return "", err
		}
	}

	return r.String(), nil
}
