package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func filePaths(args []string) ([]string, error) {
	if len(args) == 0 {
		return readDir(".")
	}

	var all []string
	for _, arg := range args {
		var err error
		var files []string

		switch {
		case arg == "./...":
			files, err = recursiveFiles(".")
		case strings.HasSuffix(arg, "/..."):
			files, err = recursivePackageFiles(arg[:len(arg)-len("/...")])
		case filepath.Ext(arg) == "":
			files, err = packageFiles(arg)
		case filepath.Ext(arg) == ".go":
			files = append(files, arg)
		}
		if err != nil {
			return nil, err
		}

		all = append(all, files...)
	}

	// Dedup files.
	dd := make(map[string]struct{})
	for _, f := range all {
		dd[f] = struct{}{}
	}

	// Convert dedup map to slice.
	var ddFiles []string
	for f := range dd {
		ddFiles = append(ddFiles, f)
	}

	return ddFiles, nil
}

func recursiveFiles(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(p string, fi os.FileInfo, err error) error {
		if fi == nil {
			return nil
		}
		if fi.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		if filepath.Ext(p) != ".go" {
			return nil
		}

		files = append(files, p)
		return nil
	})

	return files, err
}

func recursivePackageFiles(pkg string) ([]string, error) {
	pkgPath := filepath.Join(os.ExpandEnv("$GOPATH/src/"), pkg)
	return recursiveFiles(pkgPath)
}

func packageFiles(pkg string) ([]string, error) {
	pkgPath := filepath.Join(os.ExpandEnv("$GOPATH/src/"), pkg)
	return readDir(pkgPath)
}

func readDir(p string) ([]string, error) {
	fis, err := ioutil.ReadDir(p)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		if filepath.Ext(fi.Name()) != ".go" {
			continue
		}

		files = append(files, filepath.Join(p, fi.Name()))
	}

	return files, nil
}
