package disk_spider

import (
	"strings"
	"path/filepath"
	"os"
	"fmt"
)

func WalkDirToChan(dirPth string, suffixs []string, c chan *File) (files []*File, err error) {

	files = make([]*File, 0)
	i:=0
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录

		if err != nil {
			return nil
		}
		if fi == nil {
			return nil
		}
		if fi.IsDir() {
			return nil
		}

		f := &File{
			Size:     fi.Size(),
			FilePath:     filename,
			FileName: fi.Name(),
			Mode:     fi.Mode(),
		}

		for _, suffix := range suffixs {
			suffix = strings.ToUpper(suffix)
			if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
				i++
				if i>267 {
					fmt.Print("1")
				}
				f.FilePath = filename
				fmt.Println(f.FilePath)
				c <- f
				files = append(files, f)
				break
			}
		}

		return nil
	})

	return files, err
}
