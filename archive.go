package disk_spider

import (
	"bufio"
	"os"
	"errors"
	"io"
)

var (
	ErrorUninitArch = errors.New("the archive didn't init")
	ErrorProtocol   = errors.New("protocol error")
)

var fileHeadLength = 64

const (
	EncryptNone = 0
)

type Archiver struct {
	ArchiveName     string
	File            *os.File
	Writer          *bufio.Writer
	EncryptType     byte
	ProtocolVersion byte
	needRollback    bool
	CurrentTask     int
	MaxFileSize     int
}

func NewArchiver(fileName string) (*Archiver, error) {

	f, err := os.Create(fileName)

	if err != nil {
		return nil, err
	}
	w := bufio.NewWriter(f)
	arch := &Archiver{
		ArchiveName: fileName,
		File:        f,
		Writer:      w,
	}

	arch.writeFileHead()

	return arch, nil
}

func (arch *Archiver) Rollback() error {
	if arch == nil {
		return ErrorUninitArch
	}
	arch.CurrentTask --
	arch.needRollback = true

	return nil
}

func (arch *Archiver) writeFileHead() error {

	var err error

	if arch.ProtocolVersion == 0 {

		head := make([]byte, fileHeadLength)

		head[0] = arch.ProtocolVersion
		head[1] = arch.EncryptType

		arch.Writer.Write(head)

		err = arch.Writer.Flush()
		if err != nil {
			return err
		}
		return nil
	} else {

	}

	return errors.New("[file head error]")
}

func (arch *Archiver) Finish() {
	arch.File.Close()
}

func (arch *Archiver) Archive(file *File) error {

	var err error

	defer func() {
		if e := recover(); e != nil {
			arch.Rollback()
		} else {
			if err != nil {
				arch.Rollback()
			}
		}
	}()

	if arch.needRollback {
		arch.needRollback = false
	} else {
		arch.CurrentTask++
	}

	// Index
	arch.Writer.Write(bigEndianPack(uint32(arch.CurrentTask)))

	head, err := arch.generateSubHeader(file)

	if err != nil {
		return err
	}

	src, err := os.Open(file.Path)
	defer src.Close()

	if err != nil {
		return err
	}

	headLen := make([]byte, 0)
	headLen = appendBigEndian(headLen, uint32(len(head)))
	// head len
	arch.Writer.Write(headLen)
	// head
	arch.Writer.Write(head)

	bodyLen := make([]byte, 0)

	bodyLen = appendBigEndian(bodyLen, uint32(uint32(file.Size)))

	// body len
	arch.Writer.Write(bodyLen)

	n, err := io.Copy(arch.Writer, src)
	if err != nil {
		return err
	}

	if n == file.Size {
		arch.Writer.Flush()
		return nil
	}

	return nil
}

func (arch *Archiver) generateSubHeader(file *File) ([]byte, error) {

	src, err := os.Open(file.Path)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	buff := make([]byte, 0)

	if arch.EncryptType == EncryptNone {

		buff = appendBigEndian(buff, uint32(file.Size))

		fileName := []byte(file.FileName)

		buff = appendBigEndian(buff, uint32(len(fileName)))

		buff = bytesCombine(buff, fileName)

		filePath := []byte(file.Path)

		buff = appendBigEndian(buff, uint32(len(filePath)))

		buff = bytesCombine(buff, filePath)

		return buff, nil
	}

	return nil, err
}

func (arch *Archiver) CheckArchive() (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("runtime error")
		}
	}()

	offset := 0
	n := 0
	index := 0
	currentIndex := 0
	//headLen := 0
	b := make([]byte, 4)

	offset += fileHeadLength

	n, err = arch.File.ReadAt(b, int64(offset))
	if err != nil || n != 4 {
		goto FileError
	}
	index = uint32BigEndianBytesToInt(b)

	if index!=currentIndex {
		goto FileError
	}


	offset += 4
	n, err = arch.File.ReadAt(b, int64(offset))
	if err != nil || n != 4 {
		goto FileError
	}
	//headLen = uint32BigEndianBytesToInt(b)

	offset += 4



	if err != nil || n != fileHeadLength {
		goto FileError
	}

FileError:
	return errors.New("file error")

}
