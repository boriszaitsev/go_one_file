package main

import (
	"fmt"
	"os"
	"time"
)

type GoWriter struct {
	pkg           string
	imports, code []string
}

func MakeWriter() GoWriter {
	pkg := ""
	imports := make([]string, 0, 20)
	code := make([]string, 0, 4000)
	return GoWriter{pkg, imports, code}
}

func (w *GoWriter) AddGoFile(gf GoFile) {
	if len(gf.pkg) > 0 {
		w.pkg = gf.pkg
	}

	for _, fimp := range gf.imports {
		found := false
		for _, wimp := range w.imports {
			if fimp == wimp {
				found = true
				break
			}
		}

		if !found && len(fimp) != 0 {
			w.imports = append(w.imports, fimp)
		}
	}

	w.code = append(w.code, "")
	w.code = append(w.code, "//  **** "+gf.fname+" ****")
	w.code = append(w.code, "")
	w.code = append(w.code, gf.code...)
}

func (w *GoWriter) Write(path, outFname string) error {
	f, err := os.Create(path + "\\" + outFname)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintln(f, "//  "+time.Now().Format("2006-01-02 15:04"))
	fmt.Fprintln(f, w.pkg)
	fmt.Fprintln(f, "")

	fmt.Fprintln(f, "import (")
	for _, imp := range w.imports {
		fmt.Fprintln(f, "\t"+imp)
	}
	fmt.Fprintln(f, ")")
	fmt.Fprintln(f, "")

	for _, line := range w.code {
		fmt.Fprintln(f, line)
	}
	return nil
}
