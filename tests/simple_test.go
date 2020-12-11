package tests

import (
	"testing"

	rollingfile "github.com/remoteit/systemkit-logging-persisters-rollingfile"
)

func TestSimple(t *testing.T) {
	l := rollingfile.NewFileLoggerWithCustomRotationCustomNameEasy("test.log", rollingfile.Rotation{
		Count:   5,
		MaxSize: 20,
	})

	l.Info("A-1234567890-1234567890") // file 1
	l.Info("A-1234567890-1234567890") // file 2
	l.Info("A-1234567890-1234567890") // file 3
	l.Info("A-1234567890-1234567890") // file 4
	l.Info("A-1234567890-1234567890") // file 5

	l.Info("B-1234567890-1234567890") // file 1
	l.Info("B-1234567890-1234567890") // file 2
	l.Info("B-1234567890-1234567890") // file 3
	l.Info("B-1234567890-1234567890") // file 4
	l.Info("B-1234567890-1234567890") // file 5

	l.Info("C-1234567890-1234567890") // file 1
}
