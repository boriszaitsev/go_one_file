package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		pressEnterExit(err)
	}
	print("path = ", path)
	outFileName := "_out_" + filepath.Base(path) + ".go"
	print("out file name = ", outFileName)

	exePath := path
	if len(os.Args) == 2 {
		path = filepath.Join(path, os.Args[1])
	}
	dirFiles, err := MakeDirFiles(path)
	if err != nil {
		pressEnterExit(err)
		return
	}

	goList := dirFiles.ListExt("go")
	print("goList=", goList)
	writer := MakeWriter()
	for _, fname := range goList {
		if strings.Contains(fname, "_test") || strings.HasPrefix(fname, "_") {
			continue
		}

		goFile, err := ReadGoFile(path, fname)
		if err != nil {
			pressEnterExit(err)
			return
		}

		writer.AddGoFile(goFile)
	}

	if err = writer.Write(exePath, outFileName); err != nil {
		pressEnterExit(err)
		return
	}

}

// Prints error and ask user to press Enter to continue
func pressEnterExit(err error) {
	if err != nil {
		printf("ERROR: %v\n", err)
	}
	print("Нажмите Enter для выхода...")
	var input string
	fmt.Scanln(&input)
}
