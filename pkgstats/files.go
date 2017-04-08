package pkgstats

import (
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	goSrc = filepath.Join(build.Default.GOPATH, "src")
)

func packageFiles(dir string) ([]string, error) {
	var fis []os.FileInfo
	var err error
	if dir[0] == '.' {
		fis, err = ioutil.ReadDir(dir)
	} else {
		dir = filepath.Join(goSrc, dir)
		fis, err = ioutil.ReadDir(dir)
	}
	if err != nil {
		return nil, err
	}

	var files []string
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}

		switch {
		case filepath.Ext(fi.Name()) != ".go", strings.HasPrefix(fi.Name(), "_"),
			strings.HasPrefix(fi.Name(), "."):
			continue
		}

		files = append(files, filepath.Join(dir, fi.Name()))
	}

	return files, nil
}
