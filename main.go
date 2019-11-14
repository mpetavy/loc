package main

import (
	"bufio"
	"flag"
	"github.com/mpetavy/common"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type pathlist []string

type fileitem struct {
	fullname string
	relname  string
}

var paths pathlist

func init() {
	common.Init("1.0.0", "2018", "Line of code counter", "mpetavy", common.APACHE, false, nil, nil, run, 0)
	flag.Var(&paths, "i", "include directory or file")
}

func (i *pathlist) String() string {
	if i == nil {
		return ""
	}
	return strings.Join(paths, ",")
}

func (i *pathlist) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func run() error {
	files := make([]fileitem, 0)

	for _, dir := range paths {
		recursive := !strings.HasSuffix(dir, string(filepath.Separator))
		if !recursive {
			dir = dir[:len(dir)-1]
		}

		cleanDir := common.CleanPath(dir)

		err := common.WalkFilepath(cleanDir, recursive, func(file string) error {
			files = append(files, fileitem{
				fullname: file,
				relname:  file[len(filepath.Dir(cleanDir))+1:],
			})

			return nil
		})

		if common.Error(err) {
			return err
		}
	}

	cc := 0
	for _, fileItem := range files {
		ba, err := ioutil.ReadFile(fileItem.fullname)
		if common.Error(err) {
			return err
		}

		scanner := bufio.NewScanner(strings.NewReader(string(ba)))
		scanner.Split(common.ScanLinesWithLF)

		c := 0
		for scanner.Scan() {
			c++
		}

		cc += c

		common.Info("Loc %s: %d", fileItem.fullname,c)
	}

	common.Info("Sum: %d", cc)

	return nil
}

func main() {
	defer common.Done()

	common.Run([]string{"i"})
}
