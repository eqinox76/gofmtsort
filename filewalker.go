package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func Scan(output chan string, input ...string) error {
	defer close(output)
	for _, elem := range input {
		stat, err := os.Stat(elem)
		if err != nil {
			return errors.Wrapf(err, "could not stat '%s'", elem)
		}

		if !stat.IsDir() && strings.HasSuffix(elem, ".go") {
			output <- elem
			continue
		}

		err = filepath.WalkDir(elem, func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				return errors.Wrapf(err, "could not walk '%s'", path)
			}
			if !info.IsDir() && strings.HasSuffix(path, ".go") {
				output <- path
				return nil
			}
			// skip private dirs like .git
			if info.IsDir() && len(info.Name()) > 1 && strings.HasPrefix(info.Name(), ".") {
				return fs.SkipDir
			}
			return nil
		})
		if err != nil {
			return errors.Wrapf(err, "could not walk '%s'", elem)
		}
	}
	return nil
}
