package env

import (
	"log"
	"os"
	"strconv"
	"sync"
)

const (
	readerHTTPPortEnvName   = "READER_HTTP_PORT"
	readerBufferSizeEnvName = "READER_BUFFER_SIZE"
)

var (
	httpPort   = 7001
	bufferSize = 1024
	once       sync.Once
)

func GetHTTPPort() int {
	ensureInit()

	return httpPort
}

func GetBufferSize() int {
	ensureInit()

	return bufferSize
}

func ensureInit() {
	once.Do(func() {
		initHttpPort := func() {
			port := os.Getenv(readerHTTPPortEnvName)

			if port == "" {
				log.Printf("Using default server port: %d", httpPort)
				return
			}

			intPort, err := strconv.Atoi(port)

			if err != nil || intPort <= 0 {
				log.Printf("Invalid server port defined, using default one: %d", httpPort)
				return
			}

			log.Printf("Setting HTTP port: %d", intPort)

			httpPort = intPort
		}

		initBufferSize := func() {
			size := os.Getenv(readerBufferSizeEnvName)

			if size == "" {
				log.Printf("Using default buffer size: %d", bufferSize)
				return
			}

			intBufferSize, err := strconv.Atoi(size)

			if err != nil || intBufferSize <= 0 {
				log.Printf("Invalid buffer size defined, using default one: %d", bufferSize)
				return
			}

			log.Printf("Setting buffer size: %d", intBufferSize)

			bufferSize = intBufferSize

		}

		initHttpPort()
		initBufferSize()
	})
}
