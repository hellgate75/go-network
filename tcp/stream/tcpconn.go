package stream

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)
// Describe capabilities of a connection reader and buffered reader/writer/closer component
type ConnReaderWriterCloser interface {
	io.Reader
	io.WriteCloser
	// Enroll new connection for streaming purposes
	Enroll(conn net.Conn)
	// Checks if the component is running
	IsOpen() bool
	// Checks if a new stream is started and the component is reading data from the related connection
	IsReading() bool
	// Wait for a new connection is read for the first time
	Wait()
}

type rwCloser struct{
	buffer		*bytes.Buffer
	running		bool
	reading		bool
}

func (rwc *rwCloser) Read(p []byte) (n int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("ConnReaderWriterCloser.Read() - Error: %v", r))
		}
	}()
	n, err = rwc.buffer.Read(p)
	return n, err
}

func (rwc *rwCloser) Write(p []byte) (n int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("ConnReaderWriterCloser.Write() - Error: %v", r))
		}
	}()
	n, err = rwc.buffer.Write(p)
	return n, err
}
func (rwc *rwCloser) IsOpen() bool {
	return rwc.running
}

func (rwc *rwCloser) IsReading() bool {
	return rwc.reading
}

func (rwc *rwCloser) Wait() {
	for ! rwc.reading {
		time.Sleep(250 * time.Millisecond)
	}
}

func (rwc *rwCloser) Close() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("ConnReaderWriterCloser.Close() - Error: %v", r))
		}
	}()
	rwc.running = false
	rwc.buffer.Reset()
	return err
}
func (rwc *rwCloser) readFrom(conn net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(fmt.Sprintf("ConnReaderWriterCloser.readFrom() - Error: %v", r))
		}
	}()
	rwc.reading = true
	for rwc.running {
		//fmt.Println("New stream read cycle")
		var readCount int64
		data := make([]byte, 4096)
		var err error
		var n int
		n, err = conn.Read(data)
		for rwc.running && n > 0 && err == nil {
			rwc.buffer.Write(data[:n])
			readCount += int64(n)
			n, err = conn.Read(data)
		}
		if err != nil {
			time.Sleep(250 * time.Microsecond)
			continue
		}
		if readCount > 0 {
			time.Sleep(500 * time.Millisecond)
			rwc.buffer.Reset()
		}
	}
}

func (rwc *rwCloser) Enroll(conn net.Conn) {
	if conn != nil {
		if ! rwc.running {
			rwc.running = true
		}
		rwc.reading = false
		go rwc.readFrom(conn)
		time.Sleep(1 * time.Second)
	} else {
		fmt.Println("Connection is nil")
	}

}

func NewConnReaderWriterCloser() ConnReaderWriterCloser {
	return &rwCloser{
		running: false,
		reading: false,
		buffer: bytes.NewBuffer(make([]byte, 0)),
	}
}