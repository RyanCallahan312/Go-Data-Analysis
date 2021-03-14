package shared

import (
	"bufio"
	"log"
	"sync"
)

// WriteToFile writes a string to a file
func WriteToFile(data string, writer *bufio.Writer, fileLock *sync.Mutex) {

	fileLock.Lock()
	_, err := writer.WriteString(data)
	if err != nil {
		log.Panic(err)
	}

	err = writer.Flush()
	if err != nil {
		log.Panic(err)
	}
	fileLock.Unlock()
}
