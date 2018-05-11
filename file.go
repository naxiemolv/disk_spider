package disk_spider

import (
	"strings"
	"path/filepath"
	"os"
)

type File struct {
	Size int64
	Path string
	Mode os.FileMode
}

func WalkDirToChan(dirPth string, suffixs []string,c chan *File ) (files []*File, err error) {

	files = make([]*File,0)
	suffix:= suffixs[0]

	suffix = strings.ToUpper(suffix)

	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录

		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}

		f:=&File{
			Size:fi.Size(),
			Path:fi.Name(),
			Mode:fi.Mode(),
		}

		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			c<-f

			files = append(files, f)
		}

		return nil
	})

	return files, err
}