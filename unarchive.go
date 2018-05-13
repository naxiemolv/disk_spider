package disk_spider

import (
	"os"
	"bufio"
	"errors"
	"fmt"
	"io"
)

var FileHeadSize = 64
var FileError = errors.New("[archived file error]")

type UnArchiver struct {
	ArchiveName     string
	File            *os.File
	Reader          *bufio.Reader
	EncryptType     byte
	ProtocolVersion byte
	needRollback    bool
	CurrentTask     int
	MaxHeadSize     int
	MaxBodySize     int
}

func NewUnArchiver(fileName string) (*UnArchiver, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	u := &UnArchiver{
		ArchiveName: fileName,
		File:        f,
		//Reader:      r,
		MaxHeadSize:65535,
	}
	err = u.readFileHead()
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (unArch *UnArchiver) UnArchive() error {
	var err error

	headLen := 0
	bodyLen := 0
	index := 0
	n := 0
	var file *os.File
	var subFile *File
	b := make([]byte, 4)



Read:

	n, err = unArch.File.Read(b)
	if err != nil || n != 4 {
		return FileError
	}
	index = uint32BigEndianBytesToInt(b)

	if index <0 {
		return FileError
	}


	n, err = unArch.File.Read(b)
	if err != nil || n != 4 {
		return FileError
	}
	headLen = uint32BigEndianBytesToInt(b)
	subHeaderByte := make([]byte, headLen)
	if headLen > unArch.MaxHeadSize {
		return FileError
	}

	n, err = unArch.File.Read(subHeaderByte)
	if err != nil || n != headLen {
		fmt.Println("[warning:length error]")
	}

	subFile,err = deserializeBlockHead(subHeaderByte)

	if err != nil {
		return err
	}

	n, err = unArch.File.Read(b)
	bodyLen = uint32BigEndianBytesToInt(b)



	file, err = os.Create(fmt.Sprintf("[%d]-%s", index, subFile.FileName))
	if err != nil {

	} else {
		rn, err := io.CopyN(file, unArch.File, subFile.Size)
		if err != nil || rn != subFile.Size || bodyLen!=int(rn) {
			return FileError
		}
		file.Close()
	}

	goto Read

}

func (unArch *UnArchiver) parseHead(head []byte) (f *File, err error) {

	defer func() {
		if e := recover(); e != nil {
			f = nil
			err = FileError
		}
	}()

	cursor := 0

	fileSize := uint32BigEndianBytesToInt(head[0:4])
	cursor += 4

	fileNameSize := uint32BigEndianBytesToInt(head[cursor:cursor+4])

	cursor += 4
	fileName := string(head[cursor : cursor+fileNameSize])

	if len(fileName) == 0 {
		fileName = fmt.Sprintf("file-%d", unArch.CurrentTask)
	}
	cursor += fileNameSize
	pathSize := uint32BigEndianBytesToInt(head[cursor : cursor+4])

	cursor += 4
	path := string(head[cursor : cursor+pathSize])

	f = &File{
		Size:     int64(fileSize),
		FilePath:     path,
		FileName: fileName,
	}

	return f, nil
}

func (unArch *UnArchiver) readFileHead() error {
	b := make([]byte, FileHeadSize)
	n, err := unArch.File.Read(b)
	if n != FileHeadSize || err != nil {
		return FileError
	}
	unArch.ProtocolVersion = b[0]
	unArch.EncryptType = b[1]
	if unArch.ProtocolVersion != 0 {
		return FileError
	}
	return nil
}
