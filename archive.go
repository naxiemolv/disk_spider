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

	header := &FileHeader{
		Version:     0,
		EncryptType: 0,
	}
	if b, err := header.serialize(); err != nil {
		return nil, err
	} else {
		arch.Writer.Write(b)
		arch.Writer.Flush()
	}

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
	file.CurrentTask = arch.CurrentTask
	file.EncrypyType = arch.EncryptType

	// Index

	blockHead, err := file.serializeBlockHead()
	if err != nil {
		return err
	}
	arch.Writer.Write(blockHead)

	bodyLen := bigEndianPack(uint32(file.Size))

	// body len
	arch.Writer.Write(bodyLen)
	arch.Writer.Flush()
	src, err := os.Open(file.Path)
	if err != nil {
		return err
	}
	defer src.Close()

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

	if index != currentIndex {
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
