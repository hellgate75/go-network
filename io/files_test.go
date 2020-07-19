package io

import (
	"crypto/rand"
	"fmt"
	"github.com/hellgate75/go-network/testsuite"
	"io/ioutil"
	"os"
	"testing"
)

var basePath = ""

func initTest() {
	if basePath == "" {
		basePath = UniqueTempFolder("test")
		_ = os.MkdirAll(basePath, 0555)
//		fmt.Printf("Created path: %s\n", basePath)
	}
}

func tearDownTest() {
	fia, _ := ioutil.ReadDir(basePath)
	//for _, f := range fia {
	//	fmt.Printf("File: %s\n", f.Name())
	//	fmt.Printf("Size: %v\n", f.Size())
	//}
	if len(fia) > 0 {
		_ = os.RemoveAll(basePath)
		//fmt.Printf("Destroyed path: %s\n", basePath)
		basePath = ""
	}
}

func TestReadFile(t *testing.T) {
	initTest()
	defer tearDownTest()
	var path = fmt.Sprintf("%s%c%s", basePath, os.PathSeparator, "sample.txt")
	var data = make([]byte, 256)
	_, _ = rand.Read(data)
	errW := WriteFile(path, data, 0666, true)
	testsuite.AssertNil(t, "Write error must be nil", errW)
	barr, err := ReadFile(path)
//	_ = os.Remove(path)
	testsuite.AssertNil(t, "Read error must be nil", err)
	testsuite.AssertByteArraysEquals(t, "Read data array must be same than wrote data array", data, barr)
}

func TestWriteFile(t *testing.T) {
	initTest()
	defer tearDownTest()
	var path = fmt.Sprintf("%s%c%s", basePath, os.PathSeparator, "sample.txt")
	var data = make([]byte, 256)
	_, _ = rand.Read(data)
	_ = WriteFile(path, data, 0666, true)
	var fileExists = ExistsFile(path)
//	_ = os.Remove(path)
	testsuite.AssertEquals(t, fmt.Sprintf("File %s must exists", path), true, fileExists)
}