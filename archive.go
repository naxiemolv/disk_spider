package disk_spider

import (
	"bufio"
	"os"
)

type Archiver struct {
	File * os.File
	Writer *bufio.Writer
}

func NewArchiver(fileName string) (*Archiver, error){

	f,err := os.Create(fileName)

	if err != nil {
		return nil,err
	}
	w := bufio.NewWriter(f)
	arch:=&Archiver{
		File:f,
		Writer:w,
	}

	return arch,nil
}

func (arch*Archiver)Finish() {

}

func (arch*Archiver)Archive(file *File) error {
	return nil
}