package reader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Caik/go-stream-broadcast/internal/common/rest"
	"github.com/Caik/go-stream-broadcast/internal/reader/env"
)

func Serve() {
	http.HandleFunc("/upload", handleUpload)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", env.GetHTTPPort()), nil); err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	start := time.Now().UnixMilli()

	reader, err := r.MultipartReader()

	if err != nil {
		msg := fmt.Sprintf("Request is not multipart, ignoring: %v", err)

		log.Println(msg)
		rest.RespNotOk(400, w, msg, map[string]interface{}{
			"duration": time.Now().UnixMilli() - start,
		})

		return
	}

	bufferSize := env.GetBufferSize()

	for {
		part, err := reader.NextPart()

		if err != nil {
			if err == io.EOF {
				log.Println("Finished reading request")
				break
			}

			log.Printf("Error reading the request: %v", err)

			break
		}

		buf := make([]byte, bufferSize)

		log.Printf("Reading multipart named: %s", part.FormName())

		for {
			read, err := part.Read(buf)

			if err != nil && err != io.EOF {
				log.Printf("Error on reading: %v", err)
				break
			}

			if read == 0 {
				break
			}

			buf = make([]byte, bufferSize)

			if err == io.EOF {
				break
			}
		}
	}

	rest.RespOk(w, "file uploaded with success", map[string]interface{}{
		"duration": time.Now().UnixMilli() - start,
	})
}
