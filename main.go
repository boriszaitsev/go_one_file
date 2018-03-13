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
	print("base = ", filepath.Base((path)))
	// outFileName := "_out_" + filepath.Base(path) + ".go"

	dirFiles, err := MakeDirFiles(path)
	if err != nil {
		pressEnterExit(err)
		return
	}

	goList := dirFiles.ListExt("go")
	goFiles := make([]GoFile, 0, 50)
	// header := MakeHeader()
	mainIdx, constIdx := -1, -1
	for i, fname := range goList {
		if strings.Contains(fname, "_test") || strings.HasPrefix(fname, "_") {
			continue
		}
		if fname == "main.go" {
			mainIdx = i
		} else if fname == "const.go" {
			constIdx = i
		}

		goFile, err := ReadGoFile(path + "\\" + fname)
		if err != nil {
			pressEnterExit(err)
			return
		}
		// header.addPackage(goFile.pkg)
		// header.addImports(goFile.imports)

		goFiles = append(goFiles, goFile)
	}

	if constIdx != -1 {
		goFiles[0], goFiles[constIdx] = goFiles[constIdx], goFiles[0]
	}
	if mainIdx != -1 {
		swapIdx := 0
		if constIdx != 1 {
			swapIdx = 1
		}
		goFiles[swapIdx], goFiles[mainIdx] = goFiles[mainIdx], goFiles[swapIdx]
	}

	// if err = WriteOutFile(header, goFiles); err != nil {
	// 	pressEnterExit(err)
	// 	return
	// }

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
