package env

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const (
	broadcasterHTTPPortEnvName   = "BROADCASTER_HTTP_PORT"
	broadcasterBufferSizeEnvName = "BROADCASTER_BUFFER_SIZE"
	broadcasterHostList          = "BROADCASTER_HOST_LIST"
)

var (
	httpPort   = 7000
	bufferSize = 1024
	hostList   = []string{"localhost:7001", "localhost:7002"}
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

func GetHostList() []string {
	ensureInit()

	return hostList
}

func ensureInit() {
	once.Do(func() {
		initHttpPort := func() {
			port := os.Getenv(broadcasterHTTPPortEnvName)

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
			size := os.Getenv(broadcasterBufferSizeEnvName)

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

		initHostList := func() {
			list := os.Getenv(broadcasterHostList)

			if list == "" {
				log.Printf("Using default host list: %s", hostList)
				return
			}

			r := regexp.MustCompile(`\w+(:\d+)?(;\w+(:\d+)?)*`)

			if !r.MatchString(list) {
				log.Printf("Invalid host list defined, using default one: %s", hostList)
				return
			}

			log.Printf("Setting host list: %s", list)

			hostList = strings.Split(list, ";")
		}

		initHttpPort()
		initBufferSize()
		initHostList()
	})
}
