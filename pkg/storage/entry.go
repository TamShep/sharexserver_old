package storage

import (
	"io"
	"time"
)

// AuthorIdentifier identifies the author. It currently is a simple access token.
type AuthorIdentifier string

// ID can vary and is therefore mutable.
type ID interface{}

// ReadCloseSeekOpener only implements the Read/Close method if the Open method is called.
type ReadCloseSeekOpener interface {
	// Allow access via the built in interface and implement the Read, Close and Seek methods.
	io.ReadCloser
	io.Seeker

	// Open enables the Read and Close methods. It returns an error if something
	Open() error
}

// Entry represents an uploaded file and its metadata in the storage system.
type Entry struct {
	// ID is an identical token which identifies the entry.
	ID ID
	// CallReference is also an identical token but this one is used in the request uri.
	CallReference string
	// AuthorIdentifier determines the uploader information.
	Author AuthorIdentifier
	// Filename is the name of the file (contains the application name and date) which is sent with by the ShareX client.
	Filename string
	// ContentType is the MIME-Type of the uploaded file.
	ContentType string
	// UploadDate is the unix timestamp when the file was uploaded.
	UploadDate time.Time
	// ReadCloseSeekOpener allows to read the image data while controlling the reading start process.
	Reader ReadCloseSeekOpener
}
