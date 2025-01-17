package main

import (
	"os"
	"sync"
	_ "time"
	"bytes"
	"encoding/gob"
)

var rdbFile *os.File
var rdbMu sync.RWMutex = sync.RWMutex{}

func openRdb() error {
	file, err := os.OpenFile("dump.rdb", os.O_CREATE | os.O_RDWR, 0666)
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

func readRdb() error {
	b := new(bytes.Buffer)

	decoder := gob.NewDecoder(b)

	err := decoder.Decode(&store)
	if err != nil {
		return err
	}

	return nil
}

func rdb() {
	//
}
