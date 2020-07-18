package io

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/hellgate75/go-network/model/encoding"
	"gopkg.in/yaml.v2"
	"os"
)

// Unmarshal bytes and fill the given interface (pointer to structure) with given encoding type
func Unmarshal(data []byte, encodingValue encoding.Encoding, target interface{}) error {
	var err error
	defer func() {
		if r := recover(); r!= nil {
			err = errors.New(fmt.Sprintf("io.Unmashal() - Error: %v", r))
		}
	}()
	switch encodingValue {
	case encoding.EncodingJSONFormat:
		err = json.Unmarshal(data, target)
		if err != nil {
			err = errors.New(fmt.Sprintf("io.Unmashal() - Error: %v", err))
		}
	case encoding.EncodingYAMLFormat:
		err = yaml.Unmarshal(data, target)
		if err != nil {
			err = errors.New(fmt.Sprintf("io.Unmashal() - Error: %v", err))
		}
	case encoding.EncodingXMLFormat:
		err = xml.Unmarshal(data, target)
		if err != nil {
			err = errors.New(fmt.Sprintf("io.Unmashal() - Error: %v", err))
		}
	default:
		err = errors.New(fmt.Sprintf("io.Unmashal() - Error: Unknown encoding type <%v>", encodingValue))
	}
	return err
}

// Unmarshal bytes in given file path and fill the given interface (pointer to structure) with given encoding type
func UnmarshalFile(file string, encodingValue encoding.Encoding, target interface{}) error {
	var err error
	defer func() {
		if r := recover(); r!= nil {
			err = errors.New(fmt.Sprintf("io.UnmarshalFile() - Error: %v", r))
		}
	}()
	var data = make([]byte, 0)
	data, err = ReadFile(file)
	if err != nil {
		return err
	}
	err = Unmarshal(data, encodingValue, target)
	return err
}


// Marshal given interface (pointer to structure) in bytes with given encoding type
func Marshal(encodingValue encoding.Encoding, target interface{}) ([]byte, error) {
	var err error
	var data = make([]byte, 0)
	defer func() {
		if r := recover(); r!= nil {
			err = errors.New(fmt.Sprintf("io.Marshal() - Error: %v", r))
		}
	}()
	switch encodingValue {
	case encoding.EncodingJSONFormat:
		data, err = json.Marshal(target)
		if err != nil {
			err = errors.New(fmt.Sprintf("io.Marshal() - Error: %v", err))
		}
	case encoding.EncodingYAMLFormat:
		data, err = yaml.Marshal(target)
		if err != nil {
			err = errors.New(fmt.Sprintf("io.Marshal() - Error: %v", err))
		}
	case encoding.EncodingXMLFormat:
		data, err = xml.Marshal(target)
		if err != nil {
			err = errors.New(fmt.Sprintf("io.Marshal() - Error: %v", err))
		}
	default:
		err = errors.New(fmt.Sprintf("io.Marshal() - Error: Unknown encoding type <%v>", encodingValue))
	}
	return data, err
}


// Marshal given interface (pointer to structure) bytes to given file path with given encoding type
func MarshalToFile(file string, perm os.FileMode, encoding encoding.Encoding, target interface{}) error {
	var err error
	defer func() {
		if r := recover(); r!= nil {
			err = errors.New(fmt.Sprintf("io.MarshalToFile() - Error: %v", r))
		}
	}()
	var data = make([]byte, 0)
	data, err = Marshal(encoding, target)
	if err != nil {
		return err
	}
	err = WriteFile(file, data, perm, true)
	return err
}

func base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func base64Decode(str string) (string, bool) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", true
	}
	return string(data), false
}
func DecodeBase64(encoded []byte) ([]byte, error) {
	var err error
	defer func(){
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("io.DecodeBase64() - Error: %v", r))
		}
	}()
	var out []byte
	out, err = base64.StdEncoding.DecodeString(string(encoded))
	if err != nil {
		//	err = errors.New(fmt.Sprintf("io.DecodeBase64 - Decoder Read failed: %v", err))
//		return out, err
	}
	//decoder := base64.NewDecoder(base64.RawStdEncoding, bytes.NewReader(encoded))
	//readBuffer := make([]byte, base64.RawStdEncoding.DecodedLen(len(encoded)))
	//var count int
	//count, err = decoder.Read(readBuffer)
	//if count == 0 || (err != nil && err != io.EOF) {
	//	err = errors.New(fmt.Sprintf("io.DecodeBase64 - Decoder Read failed: %v", err))
	//}
	return out, err
}

func EncodeBase64(decoded []byte) ([]byte, error) {
	var err error
	defer func(){
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("io.EncodeBase64() - Error: %v", r))
		}
	}()
	var out = make([]byte, 0)
	str := base64.StdEncoding.EncodeToString(decoded)
	out = []byte(str)
	//data := bytes.NewBuffer(make([]byte, 0))
	//encoder := base64.NewEncoder(base64.RawStdEncoding, data)
	//defer func(){
	//	_ = encoder.Close()
	//}()
	//_, err = encoder.Write(decoded)
	//if err != nil {
	//	return out, errors.New(fmt.Sprintf("io.EncodeBase64 - Encoder Write failed: %v", err))
	//}
	//out = data.Bytes()
	return out, err
}