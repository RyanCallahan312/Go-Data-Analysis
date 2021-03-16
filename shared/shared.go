package shared

import (
	"bufio"
	"log"
	"sync"
)

var (
	// Writer error writer
	Writer   *bufio.Writer
	fileLock *sync.Mutex
)

// WriteToFile writes a string to a file
func WriteToFile(data string) {

	fileLock.Lock()
	_, err := Writer.WriteString(data)
	if err != nil {
		log.Panic(err)
	}

	err = Writer.Flush()
	if err != nil {
		log.Panic(err)
	}
	fileLock.Unlock()
}
