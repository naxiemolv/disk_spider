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
	File            *os.File
	Writer          *bufio.Writer
	EncryptType     byte
	ProtocolVersion byte
	CurrentTask     int
	needRollback    bool
}

func NewArchiver(fileName string) (*Archiver, error) {

	f, err := os.Create(fileName)

	if err != nil {
		return nil, err
	}
	w := bufio.NewWriter(f)
	arch := &Archiver{
		File:   f,
		Writer: w,
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

}

func (arch *Archiver) writeFileHead() error {

	if arch == nil {
		return ErrorUninitArch
	}

	var err error

	if arch.ProtocolVersion == 0 {

		head := make([]byte, 64)

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

	if arch == nil {
		return ErrorUninitArch
	}

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

	b := make([]byte, 0)
	b = appendBigEndian(b, uint32(len(head)))
	arch.Writer.Write(b)

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

		appendBigEndian(buff, uint32(file.Size))

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
