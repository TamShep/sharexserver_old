package storage_test

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/mmichaelb/sharexserver/pkg/storage"
	"io"
	"log"
	"math/rand"
	"time"
)

const idLength = 32

// TestStorage is just a test implementation of the storage.FileStorage which stores entries temporarily in memory.
type TestStorage struct {
	entries map[*storage.Entry][]byte
}

// Initialize is the implementation of the storage.FileStorage.Initialize method.
func (testStorage *TestStorage) Initialize() error {
	testStorage.entries = make(map[*storage.Entry][]byte)
	return nil
}

type closableStorageBuffer struct {
	testStorage *TestStorage
	entry       *storage.Entry
	buffer      *bytes.Buffer
}

type testReadCloseSeekOpener struct {
	data       []byte
	readSeeker *bytes.Reader
}

// Read is the implementation of the storage.Entry.Read method
func (testReadCloseSeekOpener *testReadCloseSeekOpener) Read(p []byte) (n int, err error) {
	if testReadCloseSeekOpener.readSeeker == nil {
		return 0, errors.New("the Open method has to be called first")
	}
	return testReadCloseSeekOpener.readSeeker.Read(p)
}

// Close is the implementation of the storage.Entry.Close method
func (testReadCloseSeekOpener *testReadCloseSeekOpener) Close() error {
	if testReadCloseSeekOpener.readSeeker == nil {
		return errors.New("the Open method has to be called first")
	}
	return nil
}

// Seek is the implementation of the storage.Entry.Seek method
func (testReadCloseSeekOpener *testReadCloseSeekOpener) Seek(offset int64, whence int) (int64, error) {
	if testReadCloseSeekOpener.readSeeker == nil {
		return 0, errors.New("the Open method has to be called first")
	}
	return testReadCloseSeekOpener.readSeeker.Seek(offset, whence)
}

// Open is the implementation of the storage.Entry.Open method
func (testReadCloseSeekOpener *testReadCloseSeekOpener) Open() error {
	testReadCloseSeekOpener.readSeeker = bytes.NewReader(testReadCloseSeekOpener.data)
	return nil
}

// Close is the implementation of the io.Closer interface method.
func (closableStorageBuffer *closableStorageBuffer) Close() error {
	closableStorageBuffer.testStorage.entries[closableStorageBuffer.entry] = closableStorageBuffer.buffer.Bytes()
	return nil
}

// Write is the implementation of the io.Writer interface method which calls the Write method of the buffer instance.
func (closableStorageBuffer *closableStorageBuffer) Write(p []byte) (n int, err error) {
	return closableStorageBuffer.buffer.Write(p)
}

// Store is the implementation of the storage.FileStorage.Store method.
func (testStorage *TestStorage) Store(entry *storage.Entry) (writerCloser io.WriteCloser, err error) {
idCreation:
	entry.ID = make([]byte, idLength)
	if _, err = rand.Read(entry.ID.([]byte)); err != nil {
		return
	}
	for entry := range testStorage.entries {
		if bytes.Equal(entry.ID.([]byte), entry.ID.([]byte)) {
			goto idCreation
		}
	}
	// set call reference to the basic hex string
	entry.CallReference = hex.EncodeToString(entry.ID.([]byte))
	testStorage.entries[entry] = []byte{}
	return &closableStorageBuffer{testStorage, entry, bytes.NewBuffer([]byte{})}, err
}

// Request is the implementation of the storage.FileStorage.Request method.
func (testStorage *TestStorage) Request(callReference string) (*storage.Entry, error) {
	for entry, data := range testStorage.entries {
		if entry.CallReference == callReference {
			entryCopy := *entry
			entryPointerCopy := &entryCopy
			entryPointerCopy.Reader = &testReadCloseSeekOpener{data: data}
			return entryPointerCopy, nil
		}
	}
	return nil, storage.ErrEntryNotFound
}

// Close is the implementation of the storage.FileStorage.Close method.
func (testStorage *TestStorage) Close() error {
	// no connection etc. has to be closed because the data is just in the memory
	return nil
}

// ExampleFileStorage shows the way of implementing a custom storage.FileStorage instance. It validates its functionality
// by storing and requesting the stored entry and checking the content. While normal entry storages are persistent, this
// example stores the data only in memory.
func ExampleFileStorage() {
	testBytes := []byte("Hello, this is a test!")
	var fileStorage storage.FileStorage = &TestStorage{}
	if err := fileStorage.Initialize(); err != nil {
		log.Println("There was an error while initializing the FileStorage.")
		panic(err)
	}
	testEntry := &storage.Entry{
		UploadDate:  time.Now(),
		Filename:    "testfile.png",
		ContentType: "my/mime/type",
		Author:      storage.AuthorIdentifier("a testing person"),
	}
	if writer, err := fileStorage.Store(testEntry); err != nil {
		log.Println("There was an error while storing the TestEntry.")
		panic(err)
	} else {
		if _, err = writer.Write(testBytes); err != nil {
			log.Println("There was an error while writing the data to the TestEntry writer.")
			panic(err)
		}
		if err = writer.Close(); err != nil {
			log.Println("There was an error while closing the data writer of the TestEntry writer.")
			panic(err)
		}
	}
	if requestedTestEntry, err := fileStorage.Request(testEntry.CallReference); err != nil {
		log.Println("There was an error while requesting the stored TestEntry.")
		panic(err)
	} else {
		if err := requestedTestEntry.Reader.Open(); err != nil {
			log.Println("There was an error while opening the requested TestEntry's reader.")
			panic(err)
		}
		requestTestBytes := make([]byte, len(testBytes))
		if _, err := requestedTestEntry.Reader.Read(requestTestBytes); err != nil {
			log.Println("There was an error while reading the requested TestEntry's content.")
			panic(err)
		}
		fmt.Println(bytes.Equal(testBytes, requestTestBytes))
	}
	if err := fileStorage.Close(); err != nil {
		log.Println("There was an error while closing the FileStorage.")
		panic(err)
	}
	// Output: true
}
