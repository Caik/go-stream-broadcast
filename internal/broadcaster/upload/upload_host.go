package upload

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"sync"
)

type Host struct {
	host            string
	pipeReader      *io.PipeReader
	pipeWriter      *io.PipeWriter
	multipartWriter *multipart.Writer
	formFile        io.Writer
}

func (u *Host) OpenConnection(resChan chan<- HostResponse) {
	go func() {
		res, err := http.Post(fmt.Sprintf("http://%s/upload", u.host), u.multipartWriter.FormDataContentType(), u.pipeReader)

		if err != nil {
			log.Printf("Error on uploading: %v", err)

			resChan <- HostResponse{
				Host: u.host,
				Err:  err,
			}

			return
		}

		defer func() {
			res.Body.Close()

			if err := u.pipeReader.Close(); err != nil {
				log.Printf("Error on closing reader connection: %v", err)
			}
		}()

		b, err := ioutil.ReadAll(res.Body)

		if err != nil {
			resChan <- HostResponse{
				Host: u.host,
				Err:  err,
			}

			return
		}

		data := make(map[string]interface{})

		if err = json.Unmarshal(b, &data); err != nil {
			resChan <- HostResponse{
				Host: u.host,
				Err:  err,
			}

			return
		}

		resChan <- HostResponse{
			Host:       u.host,
			StatusCode: res.StatusCode,
			Body:       data,
			Err:        nil,
		}
	}()
}

func (u *Host) NewFormFile(formName, fileName string, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()

		mw, err := u.multipartWriter.CreateFormFile(formName, fileName)

		if err != nil {
			log.Printf("Error creating form file: %v", err)

			u.formFile = nil
			return
		}

		u.formFile = mw
	}()
}

func (u *Host) Write(buf []byte, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()

		if u.formFile == nil {
			// log.Println("Form file is nil, not writing")
			return
		}

		if _, err := u.formFile.Write(buf); err != nil {
			log.Printf("Error writing to peer: %v", err)
			return
		}
	}()
}

func (u *Host) CloseConnection() {
	go func() {
		if err := u.multipartWriter.Close(); err != nil {
			log.Printf("Error on closing connection: %v", err)
		}

		if err := u.pipeWriter.Close(); err != nil {
			log.Printf("Error on closing writer connection: %v", err)
		}
	}()
}

func NewHost(host string) *Host {
	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)

	return &Host{
		host:            host,
		pipeReader:      pr,
		pipeWriter:      pw,
		multipartWriter: mw,
	}
}
