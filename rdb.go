package main

import (
	_ "bytes"
	_ "encoding/gob"
	"os"
	"sync"
	_ "time"
)

var rdbFile *os.File
var rdbMu sync.RWMutex = sync.RWMutex{}

func openRdb() error {
	file, err := os.OpenFile("dump.rdb", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	rdbFile = file

	return nil
}

func closeRdb() error {
	err := rdbFile.Close()

	return err
}
