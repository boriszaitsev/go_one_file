package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

//-----------------------------------------------------------------------------
type DirFiles struct {
	Path string
	fset map[string]struct{}
}

//-----------------------------------------------------------------------------
func MakeDirFiles(path string) (DirFiles, error) {
	dir := DirFiles{}
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return dir, errors.New("'" + path + "' does not exists")
		} else {
			return dir, err
		}
	}
	if !strings.HasSuffix(path, "\\") {
		path += "\\"
	}

	dir.Path = path
	err := dir.Refresh()
	if err != nil {
		return dir, err
	}
	return dir, nil
}

//-----------------------------------------------------------------------------
func (dir *DirFiles) Refresh() error {
	flist, err := filepath.Glob(dir.Path + "*")
	if err != nil {
		return err
	}

	dir.fset = make(map[string]struct{}, len(flist))
	for _, fullName := range flist {
		base := filepath.Base(fullName)
		dir.fset[strings.ToLower(base)] = struct{}{}
	}
	return nil
}

//-----------------------------------------------------------------------------
func (dir *DirFiles) Count() int {
	return len(dir.fset)
}

//-----------------------------------------------------------------------------
func (dir *DirFiles) ListExt(extList ...string) []string {
	list := make([]string, 0, len(dir.fset))
	for _, ext := range extList {
		ext = "." + strings.ToLower(ext)
		for fname, _ := range dir.fset {
			if filepath.Ext(fname) == ext {
				list = append(list, fname)
			}
		}
	}
	return list
}

//-----------------------------------------------------------------------------
func (dir *DirFiles) Exists(fname string) bool {
	_, found := dir.fset[strings.ToLower(fname)]
	return found
}

//-----------------------------------------------------------------------------
// if fname doesn't exists in dir, returns fname without changing
// else tryes to modify name until it become unic
func (dir *DirFiles) GenerateUnicName(fname string) string {
	for dir.Exists(fname) {
		ext := filepath.Ext(fname)
		name := fname[:len(fname)-len(ext)]
		if len(name) < 2 {
			fname = name + "_0" + ext
			continue
		}
		prefix := name[:len(name)-2]
		suffix := name[len(name)-2:]
		if suffix[0] == '_' && suffix[1] >= '0' && suffix[1] <= '8' {
			fname = prefix + "_" + string(suffix[1]+1) + ext
		} else {
			fname = name + "_0" + ext
		}
	}
	return fname
}

//-----------------------------------------------------------------------------
// renames file inside dir. Both oldName and newName are short names
// returns newName (may differ from arg)
func (dir *DirFiles) Rename(oldName, newName string) (string, error) {
	if !dir.Exists(oldName) {
		return "", errors.New("DirFiles.Rename: " +
			oldName + " does not exists")
	}
	newName = dir.GenerateUnicName(newName)
	err := os.Rename(dir.Path+oldName, dir.Path+newName)
	if err != nil {
		return "", err
	}
	delete(dir.fset, strings.ToLower(oldName))
	dir.fset[strings.ToLower(newName)] = struct{}{}
	return newName, nil
}

//-----------------------------------------------------------------------------
func (dir *DirFiles) Move(fullFromName, newName string) (string, error) {
	newName = dir.GenerateUnicName(newName)
	err := os.Rename(fullFromName, dir.Path+newName)
	if err != nil {
		return "", err
	}
	dir.fset[strings.ToLower(newName)] = struct{}{}
	return newName, nil
}

// func main() {
// 	dir, _ := MakeDirFiles("F:\\acontrol")
// 	print(dir.GenerateUnicName(os.Args[1]))
// 	newName, err := dir.Rename(os.Args[1], os.Args[1])
// 	if err != nil {
// 		print(err)
// 	} else {
// 		print(os.Args[1], " -> ", newName)
// 	}
// }
