package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

var rdbFile *os.File
var rdbMu = sync.RWMutex{}
var rdbClosed = false

func initializeRdb() {
	openRdb()
	readRdb()
	rdb()
	captureSIGINT()
}

func openRdb() error {
	file, err := os.OpenFile("dump.rdb", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	rdbFile = file

	return nil
}

func closeRdb() error {
	writeRdb()

	rdbMu.Lock()
	err := rdbFile.Close()
	if err != nil {
		return nil
	}

	defer rdbMu.Unlock()

	return nil
}

func writeRdb() {
	b := new(bytes.Buffer)
	encoder := gob.NewEncoder(b)

	storeMu.RLock()
	err := encoder.Encode(store)
	storeMu.RUnlock()

	if err != nil {
		return
	}

	rdbMu.Lock()
	rdbFile.Seek(0, 0)
	_, err = rdbFile.Write(b.Bytes())
	rdbMu.Unlock()

	if err != nil {
		return
	}

	rdbFile.Sync()
}

func readRdb() error {
	rdbMu.RLock()
	decoder := gob.NewDecoder(bufio.NewReader(rdbFile))
	rdbMu.RUnlock()

	storeMu.Lock()
	err := decoder.Decode(&store)
	storeMu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func rdb() {
	go func() {
		for {
			if rdbClosed {
				fmt.Println("rdb closed")
				return
			}

			writeRdb()

			time.Sleep(1 * time.Second)
		}
	}()
}

func captureSIGINT() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c

		rdbClosed = true

		fmt.Println("\nWriting to dump.rdb")
		closeRdb()

		os.Exit(0)
	}()
}
