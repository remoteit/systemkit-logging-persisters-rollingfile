package persisters

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	logging "github.com/remoteit/systemkit-logging"
)

// Rotation -
type Rotation struct {
	Count   int   `json:"count"`   // nr of files
	MaxSize int64 `json:"maxSize"` // max file size before move to next one
}

// NewDefaultRotation -
func NewDefaultRotation() Rotation {
	return Rotation{
		Count:   5,       // 5 files
		MaxSize: 1000000, // 1 MB
	}
}

type fileLogger struct {
	file          *os.File
	errorOccurred bool
	errorWriter   io.Writer

	rotationConfig            Rotation
	rotationInitialFileName   string
	rotationFileIndex         int
	rotationTotalWrittenBytes int64
}

// NewFileLoggerWithRotation -
func NewFileLoggerWithRotation(fileName string, errorWriter io.Writer, rotation Rotation) logging.CoreLogger {
	fl := &fileLogger{
		errorWriter:               errorWriter,
		rotationConfig:            rotation,
		rotationInitialFileName:   fileName,
		rotationFileIndex:         0,
		rotationTotalWrittenBytes: 0,
	}

	fl.closeCurrentAndCreateNext()

	return fl
}

// Log - implement `logging.CoreLogger` interface
func (thisRef *fileLogger) Log(logEntry logging.LogEntry) logging.LogEntry {
	bytesToWrite := []byte(logEntry.Message + "\n")

	if thisRef.errorOccurred && thisRef.errorWriter != nil {
		thisRef.errorWriter.Write(bytesToWrite)
	} else {
		if thisRef.rotationTotalWrittenBytes > thisRef.rotationConfig.MaxSize {
			thisRef.closeCurrentAndCreateNext()
		}

		thisRef.file.Write(bytesToWrite)
		thisRef.rotationTotalWrittenBytes += int64(len(bytesToWrite))
	}

	return logEntry
}

func (thisRef *fileLogger) closeCurrentAndCreateNext() {
	if thisRef.file != nil {
		thisRef.file.Sync()
		thisRef.file.Close()
	}

	// reset the bytes written
	thisRef.rotationTotalWrittenBytes = 0

	// generate new file name
	if thisRef.rotationFileIndex >= thisRef.rotationConfig.Count {
		thisRef.rotationFileIndex = 0
	}

	var nextFileName = thisRef.rotationInitialFileName
	if thisRef.rotationFileIndex != 0 {
		fileExtension := path.Ext(nextFileName)
		nextFileName = strings.Replace(nextFileName, fileExtension, "", 1)
		nextFileName = fmt.Sprintf("%s-%d%s", nextFileName, thisRef.rotationFileIndex, fileExtension)
	}
	thisRef.rotationFileIndex++

	// create the new file
	var err error
	if !fileOrFolderExists(nextFileName) {
		thisRef.file, err = os.Create(nextFileName)
	} else {
		thisRef.file, err = os.OpenFile(nextFileName, os.O_WRONLY|os.O_APPEND, 0660)
		if err == nil {
			fileInfo, err := thisRef.file.Stat()
			if err == nil && fileInfo != nil {
				fileSize := fileInfo.Size()
				if fileSize < thisRef.rotationConfig.MaxSize {
					thisRef.rotationTotalWrittenBytes = fileSize
				} else {
					thisRef.file.Truncate(0)
					thisRef.file.Sync()
				}
			}
		}
	}
	if err != nil && thisRef.errorWriter != nil {
		thisRef.errorWriter.Write([]byte(err.Error() + "\n"))
	}

	thisRef.errorOccurred = (err != nil)
}
