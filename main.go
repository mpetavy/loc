package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/mpetavy/common"
	"os"
	"strings"
)

type pathlist []string

var (
	paths     pathlist
	recursive *bool
)

func init() {
	common.Init("loc", "1.0.1", "", "", "2018", "Line of code counter", "mpetavy", fmt.Sprintf("https://github.com/mpetavy/%s", common.Title()), common.APACHE, nil, nil, nil, run, 0)
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
	files := make([]string, 0)

	for _, dir := range paths {
		fw, err := common.NewFilewalker(dir, *recursive, false, func(file string, f os.FileInfo) error {
			if f.IsDir() {
				return nil
			}

			common.Debug("found file: %s", file)
			files = append(files, file)

			return nil
		})
		if common.Error(err) {
			return err
		}

		err = fw.Run()
		if common.Error(err) {
			return err
		}
	}

	cc := 0
	for _, file := range files {
		ba, err := os.ReadFile(file)
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

		common.Info("%s: %d", file, c)
	}

	common.Info("Sum: %d", cc)

	return nil
}

func main() {
	common.Run([]string{"i"})
}
