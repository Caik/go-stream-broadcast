package broadcaster

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Caik/go-stream-broadcast/internal/broadcaster/env"
	"github.com/Caik/go-stream-broadcast/internal/broadcaster/upload"
	"github.com/Caik/go-stream-broadcast/internal/common/rest"
)

func Serve() {
	http.HandleFunc("/broadcast", handleBroadcast)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", env.GetHTTPPort()), nil); err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}

func handleBroadcast(w http.ResponseWriter, r *http.Request) {
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

	var wg sync.WaitGroup
	hostList := env.GetHostList()
	bufferSize := env.GetBufferSize()

	resChan := make(chan upload.HostResponse, len(hostList))
	uploadHosts := make([]*upload.Host, 0, len(hostList))

	for _, h := range hostList {
		host := upload.NewHost(h)
		uploadHosts = append(uploadHosts, host)
		host.OpenConnection(resChan)
	}

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

		wg.Add(len(uploadHosts))

		for _, host := range uploadHosts {
			host.NewFormFile(part.FormName(), part.FileName(), &wg)
		}

		wg.Wait()

		log.Printf("Reading multipart named: %s", part.FormName())

		for {
			buf := make([]byte, bufferSize)
			read, err := part.Read(buf)

			if err != nil && err != io.EOF {
				log.Printf("Error on reading: %v", err)
				break
			}

			if read == 0 {
				break
			}

			wg.Add(len(uploadHosts))

			for _, host := range uploadHosts {
				host.Write(buf, &wg)
			}

			wg.Wait()
		}
	}

	for _, host := range uploadHosts {
		host.CloseConnection()
	}

	wg.Add(len(uploadHosts))
	data := make(map[string]interface{})
	data["duration"] = time.Now().UnixMilli() - start

	go func() {
		wg.Wait()
		close(resChan)
	}()

	for res := range resChan {
		if res.Err != nil {
			data[res.Host] = map[string]string{
				"message": res.Err.Error(),
			}
		} else {
			data[res.Host] = res.Body
		}

		wg.Done()
	}

	rest.RespOk(w, "ok", data)
}
