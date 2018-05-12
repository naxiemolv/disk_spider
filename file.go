package disk_spider

import (
	"os"
)

type File struct {
	Size     int64
	Path     string
	FileName string
	Mode     os.FileMode
	File     *os.File
}
