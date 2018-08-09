// Copyright 2018 Shinichi MOTOKI. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

var (
	errUnsupportArgumentType = errors.New("unsupport argument type")
)

func (c *tmpleContext) funcMap() template.FuncMap {
	return template.FuncMap{
		"glob":            c.tfGlob,
		"includeFile":     c.tfIncludeFile,
		"includeTemplate": c.tfIncludeTemplate,
		"includeTextFile": c.tfIncludeTextFile,
	}
}

func (c *tmpleContext) tfGlob(glob string) ([]string, error) {
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
			return nil, errUnsupportArgumentType
		}
	}

	return l, nil
}

func (c *tmpleContext) tfIncludeFile(args ...interface{}) (string, error) {
	paths, err := argToStringSlice(args)
	if err != nil {
		return "", err
	}

	s := strings.Builder{}

	for _, p := range paths {
		b, err := ioutil.ReadFile(p)
		if err != nil {
			return "", err
		}
		s.Write(b)
	}

	return s.String(), nil
}

func (c *tmpleContext) tfIncludeTextFile(args ...interface{}) (string, error) {
	paths, err := argToStringSlice(args)
	if err != nil {
		return "", err
	}

	s := strings.Builder{}

	for _, p := range paths {
		b, err := ioutil.ReadFile(p)
		if err != nil {
			return "", err
		}
		s.Write(b)
		if c := b[len(b)-1]; c != '\r' && c != '\n' {
			s.WriteByte('\n')
		}
	}

	return s.String(), nil
}

func (c *tmpleContext) tfIncludeTemplate(args ...interface{}) (string, error) {
	paths, err := argToStringSlice(args)
	if err != nil {
		return "", err
	}

	b := &bytes.Buffer{}

	for _, p := range paths {
		fp, err := c.absPath(p)
		if err != nil {
			return "", err
		}

		tn, err := c.makeTemplateName(fp)
		if err != nil {
			return "", err
		}

		tb, err := ioutil.ReadFile(fp)
		if err != nil {
			return "", err
		}

		c.log.Debugf("tfIncludeTemplate: load template %s (%s)", tn, fp)

		_, err = c.tmpl.New(tn).Parse(string(tb))
		if err != nil {
			return "", err
		}

		c.fullpath[tn] = fp

		err = c.dir.run(filepath.Dir(fp), func() error {
			c.log.Debugf("tfIncludeTemplate: change working directory %s", c.dir.getcwd())
			c.log.Debugf("tfIncludeTemplate: execute template %s", tn)

			return c.tmpl.ExecuteTemplate(b, p, c.data)
		})
		if err != nil {
			return "", err
		}
	}

	return string(b.Bytes()), nil
}
