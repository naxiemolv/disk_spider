package disk_spider

import (
	"os"
)

type File struct {
	Size        int64
	FilePath    string
	FileName    string
	Mode        os.FileMode
	File        *os.File
	CurrentTask int
	EncrypyType byte
}

type FileHeader struct {
	Version     byte
	EncryptType byte
}

var fileHeadLength = 64

func (header *FileHeader) serialize() ([]byte, error) {
	if header.Version == 0 {
		head := make([]byte, fileHeadLength)
		head[0] = header.Version
		head[1] = header.EncryptType

		return head, nil
	}
	return nil, FileError
}

func parseFileHeader([]byte) (*FileHeader, error) {
	return nil, nil
}

func (f *File) serializeBlockHead() ([]byte, error) {


	buff := make([]byte, 0)

	buff = appendBigEndian(buff, uint32(f.CurrentTask))

	head := make([]byte, 0)
	if f.EncrypyType == 0 {

		head = appendBigEndian(head, uint32(f.Size))

 		fileName := []byte(f.FileName)

		head = appendBigEndian(head, uint32(len(fileName)))

		head = bytesCombine(head, fileName)

		filePath := []byte(f.FilePath)


		head = appendBigEndian(head,uint32(len(filePath)) )

		head = bytesCombine(head, filePath)


		buff = appendBigEndian(buff, uint32(len(head)))
		buff = bytesCombine(buff, head)
	}
	return buff, nil
}

func deserializeBlockHead(head []byte) (f *File, err error) {

	defer func() {
		if e := recover(); e != nil {
			f = nil
			err = FileError
		}
	}()

	cursor := 0

	fileSize := uint32BigEndianBytesToInt(head[0:4])
	cursor += 4

	fileNameSize := uint32BigEndianBytesToInt(head[cursor : cursor+4])

	cursor += 4
	fileName := string(head[cursor : cursor+fileNameSize])

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
