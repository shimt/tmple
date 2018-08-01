package main

import (
	"io"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
	"github.com/shimt/go-logif"
)

type tmpleContext struct {
	tmpl *template.Template

	base string
	dir  dirStack

	data     map[string]interface{}
	fullpath map[string]string

	log logif.LeveledLogger
}

func (c *tmpleContext) newTemplate(name string) *template.Template {
	return template.New(name).Funcs(c.funcMap())
}

func (c *tmpleContext) absPath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	p := filepath.Clean(filepath.Join(c.dir.getcwd(), path))

	return p, nil
}

func (c *tmpleContext) makeTemplateName(path string) (string, error) {
	p, err := c.absPath(path)
	if err != nil {
		return "", errors.Wrap(err, "template name error")
	}

	p, err = filepath.Rel(c.base, p)
	if err != nil {
		return "", errors.Wrap(err, "template name error")
	}

	return p, nil
}

func (c *tmpleContext) Execute(out io.Writer) (err error) {
	return c.dir.run(c.base, func() error {
		return c.tmpl.Execute(out, c.data)
	})
}
