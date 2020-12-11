package persisters

import (
	"fmt"
	"os"

	logging "github.com/remoteit/systemkit-logging"
)

func fileOrFolderExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

type emptyWritter struct{}

func (thisRef emptyWritter) Write(p []byte) (n int, err error) {
	return 0, nil
}

func NewFileLoggerWithCustomRotationCustomName(fileName string, rotation Rotation) logging.CoreLogger {
	return NewFileLoggerWithRotation(fileName, &emptyWritter{}, rotation)
}

func NewFileLoggerWithDefaultRotationCustomName(fileName string) logging.CoreLogger {
	return NewFileLoggerWithCustomRotationCustomName(fileName, NewDefaultRotation())
}

func NewFileLoggerWithDefaultRotationDefaultName() logging.CoreLogger {
	return NewFileLoggerWithDefaultRotationCustomName(fmt.Sprintf("%s.log", os.Args[0]))
}

func NewFileLoggerWithDefaultRotationCustomNameEasy(fileName string) logging.Logger {
	return logging.NewLoggerImplementation(NewFileLoggerWithDefaultRotationCustomName(fileName))
}

func NewFileLoggerWithDefaultRotationDefaultNameEasy() logging.Logger {
	return NewFileLoggerWithDefaultRotationCustomNameEasy(fmt.Sprintf("%s.log", os.Args[0]))
}

func NewFileLoggerWithCustomRotationCustomNameEasy(fileName string, rotation Rotation) logging.Logger {
	return logging.NewLoggerImplementation(NewFileLoggerWithRotation(fileName, &emptyWritter{}, rotation))
}
