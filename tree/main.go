package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FileProps struct {
	Name         string
	RelativePath string
	Size         int64
	isDir        bool
	Children     []FileProps
}

func (fp *FileProps) sortChildrenByName() {
	sort.Slice(fp.Children, func(i, j int) bool {
		return fp.Children[i].Name < fp.Children[j].Name
	})
}

func addChild(path string, relativeDir string, parent *FileProps) {
	files, err := os.ReadDir(path)
	if err != nil {
		panic("Could not read directory tree")
	}
	parent.Children = make([]FileProps, len(files))

	for i, file := range files {
		parent.Children[i].Name = file.Name()

		if file.IsDir() {
			parent.Children[i].isDir = true
			parent.Children[i].RelativePath = filepath.Join(relativeDir, parent.Children[i].Name)
			addChild(filepath.Join(path, parent.Children[i].Name), parent.Children[i].RelativePath, &parent.Children[i])
		} else {
			parent.Children[i].isDir = false
			parent.Children[i].RelativePath = relativeDir
			if i == len(files)-1 {
				parent.sortChildrenByName()
			}
			//fileInfo, err2 := file.Info()
			//if err2 != nil {
			//	panic("Could not read file info")
			//}
			//fileSize := fileInfo.Size()
			//var fileSizeString string
			//if fileSize == 0 {
			//	fileSizeString = "empty"
			//} else {
			//	fileSizeString = strconv.FormatInt(fileSize, 10) + "b"
			//}
			//
			//fmt.Printf("%s (%s)\n", file.Name(), fileSizeString)
		}
	}
}

func readDir(path string) {
	result := FileProps{
		Name:         path,
		RelativePath: path + "/",
		Size:         0,
		isDir:        true,
	}
	addChild(path, result.RelativePath, &result)
	drawTree(result, 0, false)

	//fmt.Println(result)
}

func drawTree(dir FileProps, gap int, lastParent bool) {
	for index, child := range dir.Children {
		if child.isDir {
			fmt.Printf("%s├───%s\n", strings.Repeat("│\t", gap), child.Name)
			drawTree(child, gap+1, index == len(dir.Children)-1)
		} else {
			//var divider string
			//if !lastParent {
			//	divider = "│"
			//}
			stringToRepeat := "│" + "\t"

			if index == len(dir.Children)-1 {
				fmt.Printf("%s└───%s\n", strings.Repeat(stringToRepeat, gap), child.Name)
			} else {
				fmt.Printf("%s├───%s\n", strings.Repeat(stringToRepeat, gap), child.Name)
			}
		}
	}
}

func dirTree(output io.Writer, path string, printFiles bool) error {
	readDir("testdata")

	return nil
}

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
