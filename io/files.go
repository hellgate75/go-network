package io

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
)

var (
	DefaultFilePerm   os.FileMode = 0664
	DefaultFolderPerm os.FileMode = 0664
)
const BufferSize = 1024
func ReadFile(path string) ([]byte, error) {
	var out = make([]byte, 0)
	var err error
	defer func() {
		if r := recover(); r!= nil {
			err = errors.New(fmt.Sprintf("io.ReadFile() - Error: %v", r))
		}
	}()
	var fstat os.FileInfo
	if fstat, err = os.Stat(path); err != nil {
		return out, err
	}
	var file *os.File
	if file, err = os.Open(path); err != nil {
		return out, err
	}
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
		err = file.Close()
	}()
	var size = fstat.Size()
	var b = make([]byte, BufferSize)
	var cycles = int64(math.Ceil(float64(size) / BufferSize))
	var n int
	var bBuff = bytes.NewBuffer([]byte{})
	for i := int64(0); i < cycles; i++ {
		n, err = file.Read(b)
		if n <= BufferSize {
			bBuff.Write(b[:n])

		} else {
			bBuff.Write(b)
		}
	}
	out = append(out, bBuff.Bytes()...)
//	bBuff.Reset()
	return out, err
}

func WriteFile(path string, data []byte, perm os.FileMode, override bool) error {
	var err error
	defer func() {
		if r := recover(); r!= nil {
			err = errors.New(fmt.Sprintf("io.WriteFile() - Error: %v", r))
		}
	}()
	if _, err = os.Stat(path); err == nil && ! override {
		return errors.New(fmt.Sprintf("File alredy exists and no override policy adopted"))
	}
	var f *os.File
	f, err = os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
		_ = f.Sync()
		_ = f.Close()
	}()
	_ = f.Chmod(perm)
	var n int
	n, err = f.Write(data)
	if n != len(data) {
		return errors.New(fmt.Sprintf("Expected written <%v> bytes but wrote <%v>", len(data), n))
	}
	return err
}

func ExistsFile(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}


func IsFolder(path string) bool {
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return true
	}
	return false
}
func CreateFolders(path string, perm os.FileMode) error {
	var err error
	if ExistsFile(path) {
		return errors.New(fmt.Sprintf("File or folder %s already exists", path))
	}
	if IsFolder(path) {
		err = os.MkdirAll(path, perm)
	} else {
		fld, _ := filepath.Split(path)
		err = os.MkdirAll(fld, perm)
	}
	return err
}
