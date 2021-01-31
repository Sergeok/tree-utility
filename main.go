package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	//"strings"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func addFileSize(file os.FileInfo) string {
	if file.Size() == 0 {
		return file.Name() + " (empty)"
	} else {
		return file.Name() + " (" + strconv.FormatInt(file.Size(), 10) + "b)"
	}
}

type fileOrDir struct {
	name string
	isDir bool
}

func dirParser(out io.Writer, path string, printFiles bool, prefix string) {
	//list of current directory files
	var fileList []fileOrDir

	//go to current path
	err := os.Chdir(path)
	check(err)

	//getting a list of all files in the current directory
	wd, err := os.Getwd()
	check(err)
	files, err := ioutil.ReadDir(wd)

	//choosing files(names and types) by the condition
	for _, file := range files {
		if file.IsDir() {
			fileList = append(fileList, fileOrDir{name: file.Name(), isDir: file.IsDir()})
		} else if printFiles {
			fileList = append(fileList, fileOrDir{name: addFileSize(file), isDir: file.IsDir()})
		}
	}

	//file printing and recursive parsing
	if len(fileList) > 0 {
		for _, file := range fileList[:len(fileList)-1] {
			_, err = fmt.Fprint(out, prefix, "├───", file.name, "\n")
			check(err)
			if file.isDir {
				dirParser(out, file.name, printFiles, prefix + "│\t")
			}
		}

		_, err = fmt.Fprint(out, prefix, "└───", fileList[len(fileList)-1].name, "\n")
		check(err)
		if fileList[len(fileList)-1].isDir {
			dirParser(out, fileList[len(fileList)-1].name, printFiles, prefix + "\t")
		}
	}

	//go to parent directory
	err = os.Chdir(filepath.Dir(wd))
	check(err)
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	dirParser(out, path, printFiles, "")

	//error plug (panic is instantaneous and causes the program to terminate)
	return nil
}

////alternative file printing and recursive parsing part
//for i, file := range fileList {
//	if i != len(fileList)-1 {
//		fmt.Fprint(out, prefix, "├───", file.name, "\n")
//		if file.isDir {
//			dirParser(out, file.name, printFiles, prefix + "│\t")
//		}
//	} else {
//		fmt.Fprint(out, prefix, "└───", file.name, "\n")
//		if file.isDir {
//			dirParser(out, file.name, printFiles, prefix + "\t")
//		}
//	}
//}
