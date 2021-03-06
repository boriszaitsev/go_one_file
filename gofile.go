package main

import (
	"bufio"
	"os"
	"path/filepath"
	ss "strings"
)

type GoFile struct {
	fname, pkg    string
	imports, code []string
}

func ReadGoFile(path, fname string) (GoFile, error) {
	pkg := ""
	imports := make([]string, 0, 10)
	code := make([]string, 0, 200)
	file, err := os.Open(path + string(filepath.Separator) + fname)
	if err != nil {
		return GoFile{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inImport := false
	for scanner.Scan() {
		line := scanner.Text()
		trm := ss.TrimSpace(line)

		if len(trm) == 0 && ss.HasPrefix(trm, "//") {
			continue
		}

		if ss.HasPrefix(trm, "package ") {
			pkg = line
			continue
		}

		if inImport {
			imp := takeImport(trm)
			if len(imp) > 0 {
				imports = append(imports, imp)
			}
			if ss.HasSuffix(trm, ")") {
				inImport = false
			}
			continue
		}

		if ss.HasPrefix(trm, "import ") {
			if ss.Contains(trm, "(") {
				imp := takeImport(trm)
				if len(imp) > 0 {
					imports = append(imports, imp)
				}
				inImport = true
			} else {
				imports = append(imports, takeImport(trm))
			}
			continue
		}

		code = append(code, line)
	}

	if err := scanner.Err(); err != nil {
		return GoFile{}, err
	}
	print(fname)
	return GoFile{fname, pkg, imports, code}, nil
}

func takeImport(s string) string {
	fi := ss.Index(s, "\"")
	li := ss.LastIndex(s, "\"") + 1
	if fi == -1 {
		return ""
	}
	return s[fi:li]
}
