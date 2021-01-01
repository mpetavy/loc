package main

import (
	"bufio"
	"flag"
	"fmt"
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

var (
	paths     pathlist
	recursive *bool
)

func init() {
	common.Init(false, "1.0.1", "", "2018", "Line of code counter", "mpetavy", fmt.Sprintf("https://github.com/mpetavy/%s", common.Title()), common.APACHE, nil, nil, run, 0)
	flag.Var(&paths, "i", "include directory or file")
	recursive = flag.Bool("r", false, "recursive file scanning")
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
		cleanDir := common.CleanPath(dir)
		mask := ""

		common.Debug("dir: %s", cleanDir)
		common.Debug("mask: %s", mask)

		if common.ContainsWildcard(cleanDir) {
			mask = filepath.Base(cleanDir)
			cleanDir = filepath.Dir(cleanDir)
		}

		err := common.WalkFilepath(cleanDir, *recursive, false, func(file string) error {
			var err error

			b := mask == ""
			if !b {
				b, err = common.EqualWildcards(filepath.Base(file), mask)
				if common.Error(err) {
					return err
				}
			}

			if b {
				common.Debug("found file: %s", file)
				files = append(files, fileitem{
					fullname: file,
					relname:  file[len(filepath.Dir(cleanDir))+1:],
				})
			}

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

		common.Info("Loc %s: %d", fileItem.fullname, c)
	}

	common.Info("Sum: %d", cc)

	return nil
}

func main() {
	defer common.Done()

	common.Run([]string{"i"})
}
