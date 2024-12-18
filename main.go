package main

import (
	"bufio"
	"embed"
	"flag"
	"github.com/mpetavy/common"
	"os"
	"strings"
)

var (
	paths     common.MultiValueFlag
	recursive *bool
)

//go:embed go.mod
var resources embed.FS

func init() {
	common.Init("", "", "", "", "Line of code counter", "", "", "", &resources, nil, nil, run, 0)
	flag.Var(&paths, "f", "include directory or file")
	recursive = flag.Bool("r", false, "recursive file scanning")
}

func run() error {
	files := make([]string, 0)

	for _, dir := range paths {
		err := common.WalkFiles(dir, *recursive, false, func(file string, f os.FileInfo) error {
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
	common.Run([]string{"f"})
}
