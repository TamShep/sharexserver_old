package storages

import (
	"bytes"
	"errors"
	"github.com/mmichaelb/sharexserver/pkg/storage"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	// Status constants for entries
	statusWaiting = iota
	statusActivated
	statusFailed
	// Data to generate new call references
	callReferenceChars  = "abcdefghijklmnopqrstuvxyzABCDEFGHIJKLMNOPQRSTUVXYZ1234567890"
	callReferenceLength = 6
	// MongoDB key names
	iDField            = "_id"
	statusField        = "status"
	callReferenceField = "call_reference"
	authorField        = "author"
	filenameField      = "filename"
	contentTypeField   = "content_type"
	uploadDateField    = "upload_date"
)

// MongoStorage is the FileStorage implementation for the Database MongoDB in combination with the file data stored in
// standard system files.
type MongoStorage struct {
	// DialInfo contains the information used to connect to the MongoDB server.
	DialInfo *mgo.DialInfo
	// DatabaseName and CollectionName define the MongoDB specific names used to store data.
	DatabaseName, CollectionName string
	// DataFolder is the folder where uploaded files are stored in. This can be an absolute or a relative path. It has
	// to end with a slash ("/").
	DataFolder string
	// internal values
	session *mgo.Session
}

// StatusChangeWriteCloser is an extend implementation if io.FileWriter to update the database entry on close.
type StatusChangeWriteCloser struct {
	// Collection is used to update the database entry.
	Collection *mgo.Collection
	// ID is used to find and update the entry.
	ID bson.ObjectId
	// Real writer which is used to process the data.
	RealWriteCloser io.WriteCloser
}

// Write just calls the real writer to process the data.
func (writeCloser *StatusChangeWriteCloser) Write(p []byte) (int, error) {
	return writeCloser.RealWriteCloser.Write(p)
}

// Close is the extended function which also updates the database entry.
func (writeCloser *StatusChangeWriteCloser) Close() (err error) {
	var updatedStatus int
	if err = writeCloser.RealWriteCloser.Close(); err != nil {
		// set status to failed because an error occurred
		updatedStatus = statusFailed
	} else {
		// set status to activated because the data was successfully written
		updatedStatus = statusActivated
	}
	// update database entry
	if mongoErr := writeCloser.Collection.UpdateId(writeCloser.ID, bson.M{"$set": bson.M{statusField: updatedStatus}}); err != nil {
		log.Printf("An error occurred while updating the status of %v, %T: %+v",
			strconv.Quote(writeCloser.ID.String()), mongoErr, mongoErr)
	}
	return
}

func (mongoStorage *MongoStorage) Initialize() (err error) {
	// create folder for stored files
	if err = os.MkdirAll(mongoStorage.DataFolder, os.ModePerm); err != nil {
		return
	}
	// connect to MongoDB server
	mongoStorage.session, err = mgo.DialWithInfo(mongoStorage.DialInfo)
	return
}

func (mongoStorage *MongoStorage) Store(entry *storage.Entry) (writer io.WriteCloser, err error) {
	// use the provided collection to store the data in
	collection := mongoStorage.session.DB(mongoStorage.DatabaseName).C(mongoStorage.CollectionName)
randomCreation:
	// create a new random ID and call reference
	objectId := bson.NewObjectId()
	entry.ID = objectId
	entry.CallReference = mongoStorage.newCallReference()
	// insert the file details into the collection
	if err = collection.Insert(
		bson.M{
			iDField:            entry.ID,
			statusField:        statusWaiting, // set entry to waiting because the file data is not stored yet
			callReferenceField: entry.CallReference,
			authorField:        entry.Author,
			filenameField:      entry.Filename,
			contentTypeField:   entry.ContentType,
			uploadDateField:    entry.UploadDate,
		},
	); err != nil {
		if lastErr, ok := err.(*mgo.LastError); ok && lastErr.Code == 11000 {
			// duplicate key error
			goto randomCreation
		} else {
			// just return the raw error if something different happened
			return
		}
	}
	// open file and return a StatusChangeWriteCloser
	writer, err = os.Create(mongoStorage.DataFolder + objectId.Hex())
	if err != nil {
		return nil, err
	}
	// wrap the writer into an instance of the StatusChangeWriteCloser to change the status after completing the upload
	return &StatusChangeWriteCloser{
		Collection:      collection,
		ID:              objectId,
		RealWriteCloser: writer,
	}, nil
}

// method which randomly creates a new call reference
func (mongoStorage *MongoStorage) newCallReference() string {
	buf := bytes.NewBuffer([]byte{})
	for i := 0; i < callReferenceLength; i++ {
		randomIndex := rand.Intn(len(callReferenceChars))
		buf.WriteString(callReferenceChars[randomIndex : randomIndex+1])
	}
	return buf.String()
}

// FileBasedReadCloseSeekOpener is the file based implementation of the ReadCloseSeekOpener which opens a file when
// calling the Open method
type FileBasedReadCloseSeekOpener struct {
	// Filepath is used to open the file when calling the Open method
	Filepath string
	// internal values
	file *os.File
}

// Read simply just calls the real Read method and can not be called until the Open method was.
func (fileBasedReadCloseSeekOpener *FileBasedReadCloseSeekOpener) Read(p []byte) (n int, err error) {
	if fileBasedReadCloseSeekOpener.file == nil {
		return -1, errors.New("the Open method has not been called yet")
	}
	return fileBasedReadCloseSeekOpener.file.Read(p)
}

// Close simply just calls the real Close method and can not be called until the Open method was.
func (fileBasedReadCloseSeekOpener *FileBasedReadCloseSeekOpener) Close() (err error) {
	if fileBasedReadCloseSeekOpener.file == nil {
		return errors.New("the Open method has not been called yet")
	}
	return fileBasedReadCloseSeekOpener.file.Close()
}

// Seek simply just calls the real Seek method and can not be called until the Open method was.
func (fileBasedReadCloseSeekOpener *FileBasedReadCloseSeekOpener) Seek(offset int64, whence int) (int64, error) {
	if fileBasedReadCloseSeekOpener.file == nil {
		return -1, errors.New("the Open method has not been called yet")
	}
	return fileBasedReadCloseSeekOpener.file.Seek(offset, whence)
}

// Open opens the file located at the given filepath
func (fileBasedReadCloseSeekOpener *FileBasedReadCloseSeekOpener) Open() (err error) {
	fileBasedReadCloseSeekOpener.file, err = os.Open(fileBasedReadCloseSeekOpener.Filepath)
	return
}

func (mongoStorage *MongoStorage) Request(callReference string) (*storage.Entry, error) {
	collection := mongoStorage.session.DB(mongoStorage.DatabaseName).C(mongoStorage.CollectionName)
	// read result to a simple bson map
	result := &bson.M{}
	// find the entry by its call reference
	if err := collection.Find(bson.M{callReferenceField: callReference}).One(result); err == mgo.ErrNotFound {
		// return error that entry was not found
		return nil, storage.ErrEntryNotFound
	} else if err != nil {
		// return unwrapped error because something gone horrifically wrong
		return nil, err
	}
	// set all entry values except for the reader
	entry := &storage.Entry{
		ID:            storage.ID((*result)[iDField]),
		CallReference: (*result)[callReferenceField].(string),
		Author:        storage.AuthorIdentifier((*result)[authorField].(string)),
		Filename:      (*result)[filenameField].(string),
		ContentType:   (*result)[contentTypeField].(string),
		UploadDate:    (*result)[uploadDateField].(time.Time),
	}
	// initiate file based ReadCloseSeekOpener
	entry.Reader = &FileBasedReadCloseSeekOpener{
		Filepath: mongoStorage.DataFolder + (*result)[iDField].(bson.ObjectId).Hex(),
	}
	return entry, nil
}

func (mongoStorage *MongoStorage) Close() error {
	// close connection to MongoDB server
	mongoStorage.session.Close()
	return nil
}
