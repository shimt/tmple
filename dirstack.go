package main

import (
	"os"

	"github.com/pkg/errors"
)

var (
	errDirectoryStackEmpty = errors.New("directory stack empty")
)

type dirStack struct {
	stack []string
	cwd   string
}

func (s *dirStack) init() {
	s.stack = make([]string, 0, 16)
}

func (s *dirStack) pushd(path string) (err error) {
	if s.stack == nil {
		s.init()
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.Chdir(path)
	if err != nil {
		return err
	}

	s.stack = append(s.stack, cwd)
	s.cwd = path

	return nil
}

func (s *dirStack) popd() (err error) {
	if s.stack == nil {
		s.init()
	}

	l := len(s.stack)

	if l == 0 {
		return errDirectoryStackEmpty
	}

	path := s.stack[l-1]

	err = os.Chdir(path)
	if err != nil {
		return err
	}

	s.stack = s.stack[0 : l-1]
	s.cwd = path

	return nil
}

func (s *dirStack) run(path string, f func() error) (err error) {
	if err = s.pushd(path); err != nil {
		return err
	}

	defer func() {
		if e := s.popd(); e != nil {
			err = e
		}
	}()

	if err = f(); err != nil {
		return err
	}

	return nil
}

func (s *dirStack) getcwd() (path string) {
	return s.cwd
}
