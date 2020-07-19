package io

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hellgate75/go-cron/io"
	"github.com/hellgate75/go-network/model"
	"github.com/hellgate75/go-network/testsuite"
	"io/ioutil"
	"os"
	"testing"
)

var testJsonDataBytes = []byte{123,34,110,97,109,101,34,58,34,70,97,98,114,105,122,105,111,34,44,34,115,117,114,110,97,109,101,34,58,34,84,111,114,101,108,108,105,34,44,34,97,103,101,34,58,52,53,125}
var testYamlDataBytes = []byte{110,97,109,101,58,32,70,97,98,114,105,122,105,111,10,115,117,114,110,97,109,101,58,32,84,111,114,101,108,108,105,10,97,103,101,58,32,52,53,10}
var testXmlDataBytes =  []byte{60,115,97,109,112,108,101,88,77,76,62,60,110,97,109,101,62,70,97,98,114,105,122,105,111,60,47,110,97,109,101,62,60,115,117,114,110,97,109,101,62,84,111,114,101,108,108,105,60,47,115,117,114,110,97,109,101,62,60,97,103,101,62,52,53,60,47,97,103,101,62,60,47,115,97,109,112,108,101,88,77,76,62}
var testBase64EncodedData = []byte{86,71,104,112,99,121,66,112,99,121,66,104,73,72,82,108,99,51,81,61}
var testBase64DecodedData = []byte{84,104,105,115,32,105,115,32,97,32,116,101,115,116}

func TestMarshalJson(t *testing.T) {
	var tStruct = struct {
		Name		string `json:"name,omitempty"`
		Surname		string `json:"surname,omitempty"`
		Age			int 	`json:"age,omitempty"`
	}{
		"Fabrizio",
		"Torelli",
		45,
	}
	bytes, err := Marshal(model.EncodingJSONFormat, &tStruct)
	testsuite.AssertNil(t, "Marshal operation error must be nil", err)
	testsuite.AssertByteArraysEquals(t, "Json data array must be same", testJsonDataBytes, bytes)

}


func TestUnmarshalJson(t *testing.T) {
	var eStruct = struct {
		Name		string `json:"name,omitempty"`
		Surname		string `json:"surname,omitempty"`
		Age			int 	`json:"age,omitempty"`
	}{
		"Fabrizio",
		"Torelli",
		45,
	}
	var tStruct = struct {
		Name		string `json:"name,omitempty"`
		Surname		string `json:"surname,omitempty"`
		Age			int 	`json:"age,omitempty"`
	}{}
	err := Unmarshal(testJsonDataBytes, model.EncodingJSONFormat, &tStruct)
	testsuite.AssertNil(t, "Unmarshal operation error must be nil", err)
	testsuite.AssertEquals(t, "Structure must be same: name", eStruct.Name, tStruct.Name)
	testsuite.AssertEquals(t, "Structure must be same: surname", eStruct.Surname, tStruct.Surname)
	testsuite.AssertEquals(t, "Structure must be same: age", eStruct.Age, tStruct.Age)
}

func TestMarshalYaml(t *testing.T) {
	var tStruct = struct {
		Name		string `yaml:"name,omitempty"`
		Surname		string `yaml:"surname,omitempty"`
		Age			int 	`yaml:"age,omitempty"`
	}{
		"Fabrizio",
		"Torelli",
		45,
	}
	bytes, err := Marshal(model.EncodingYAMLFormat, &tStruct)
	testsuite.AssertNil(t, "Marshal operation error must be nil", err)
	testsuite.AssertByteArraysEquals(t, "Yaml data array must be same", testYamlDataBytes, bytes)
}

func TestUnmarshalYaml(t *testing.T) {
	var eStruct = struct {
		Name		string `yaml:"name,omitempty"`
		Surname		string `yaml:"surname,omitempty"`
		Age			int 	`yaml:"age,omitempty"`
	}{
		"Fabrizio",
		"Torelli",
		45,
	}
	var tStruct = struct {
		Name		string `yaml:"name,omitempty"`
		Surname		string `yaml:"surname,omitempty"`
		Age			int 	`yaml:"age,omitempty"`
	}{}
	err := Unmarshal(testYamlDataBytes, model.EncodingYAMLFormat, &tStruct)
	testsuite.AssertNil(t, "Unmarshal operation error must be nil", err)
	testsuite.AssertEquals(t, "Structure must be same: name", eStruct.Name, tStruct.Name)
	testsuite.AssertEquals(t, "Structure must be same: surname", eStruct.Surname, tStruct.Surname)
	testsuite.AssertEquals(t, "Structure must be same: age", eStruct.Age, tStruct.Age)
}

type sampleXML struct {
	Name		string `xml:"name,omitempty"`
	Surname		string `xml:"surname,omitempty"`
	Age			int 	`xml:"age,omitempty"`
}

func TestMarshalXml(t *testing.T) {
	var tStruct = sampleXML{
		"Fabrizio",
		"Torelli",
		45,
	}
	bytes, err := Marshal(model.EncodingXMLFormat, &tStruct)
	testsuite.AssertNil(t, "Marshal operation error must be nil", err)
	testsuite.AssertByteArraysEquals(t, "Xml data array must be same", testXmlDataBytes, bytes)
}

func TestUnmarshalXml(t *testing.T) {
	var eStruct = sampleXML{
		"Fabrizio",
		"Torelli",
		45,
	}
	var tStruct = sampleXML{}
	err := Unmarshal(testXmlDataBytes, model.EncodingXMLFormat, &tStruct)
	testsuite.AssertNil(t, "Unmarshal operation error must be nil", err)
	testsuite.AssertEquals(t, "Structure must be same: name", eStruct.Name, tStruct.Name)
	testsuite.AssertEquals(t, "Structure must be same: surname", eStruct.Surname, tStruct.Surname)
	testsuite.AssertEquals(t, "Structure must be same: age", eStruct.Age, tStruct.Age)
}

func generateFilePath() string {
	dir := fmt.Sprintf("%s", io.HomeFolder())
	if ! ExistsFile(dir) {
		_ = CreateFolders(dir, 0777)
	}
	uuidStr := uuid.New().String()
	return fmt.Sprintf("%s%c%s.json", dir, os.PathSeparator, uuidStr)
}

func TestMarshalFile(t *testing.T) {
	path := generateFilePath()
	var tStruct = struct {
		Name		string `json:"name,omitempty"`
		Surname		string `json:"surname,omitempty"`
		Age			int 	`json:"age,omitempty"`
	}{
		"Fabrizio",
		"Torelli",
		45,
	}
	err := MarshalToFile(path, 0777, model.EncodingJSONFormat, &tStruct)
	testsuite.AssertNil(t, "Marshal operation error must be nil", err)
	defer func() {
		_ = os.Remove(path)
	}()
	bytes, err := ioutil.ReadFile(path)
	testsuite.AssertByteArraysEquals(t, "File data byte array must be same as test data", testJsonDataBytes, bytes)
}

func TestUnMarshalFile(t *testing.T) {
	path := generateFilePath()
	var eStruct = struct {
		Name		string `json:"name,omitempty"`
		Surname		string `json:"surname,omitempty"`
		Age			int 	`json:"age,omitempty"`
	}{
		"Fabrizio",
		"Torelli",
		45,
	}
	var tStruct = struct {
		Name		string `json:"name,omitempty"`
		Surname		string `json:"surname,omitempty"`
		Age			int 	`json:"age,omitempty"`
	}{}
	err := ioutil.WriteFile(path, testJsonDataBytes, 0777)
	testsuite.AssertNil(t, "File bytes write operation error must be nil", err)
	err = UnmarshalFile(path, model.EncodingJSONFormat, &tStruct)
	testsuite.AssertNil(t, "Unmarshal operation error must be nil", err)
	defer func() {
		_ = os.Remove(path)
	}()
	testsuite.AssertNil(t, "Unmarshal operation error must be nil", err)
	testsuite.AssertEquals(t, "Structure must be same: name", eStruct.Name, tStruct.Name)
	testsuite.AssertEquals(t, "Structure must be same: surname", eStruct.Surname, tStruct.Surname)
	testsuite.AssertEquals(t, "Structure must be same: age", eStruct.Age, tStruct.Age)
}

func TestEncodeBase64(t *testing.T) {
	bytes, err := EncodeBase64(testBase64DecodedData)
	testsuite.AssertNil(t, "Marshal operation error must be nil", err)
	testsuite.AssertByteArraysEquals(t, "Base data array must be same", testBase64EncodedData, bytes)
}

func TestDecodeBase64(t *testing.T) {
	bytes, err := DecodeBase64(testBase64EncodedData)
	testsuite.AssertNil(t, "Marshal operation error must be nil", err)
	testsuite.AssertByteArraysEquals(t, "Base data array must be same", testBase64DecodedData, bytes)
}