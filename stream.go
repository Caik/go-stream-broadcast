package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

func main() {
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/copy/first", handleFirst)
	http.HandleFunc("/copy/second", handleSecond)
	log.Fatal(http.ListenAndServe(":7000", nil))
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()

	if err != nil {
		log.Printf("Is not multipart: %v\n", err)
		writeNotOk(w)

		return
	}

	pipeR1, pipeW1 := io.Pipe()
	pipeR2, pipeW2 := io.Pipe()
	m1 := multipart.NewWriter(pipeW1)
	m2 := multipart.NewWriter(pipeW2)

	go func() {
		res, err := http.Post("http://localhost:7000/copy/first", m1.FormDataContentType(), pipeR1)

		if err != nil {
			log.Printf("Error on POST first: %v\n", err)
		}

		fmt.Printf("StatusCode first: %d\n", res.StatusCode)
	}()

	go func() {
		res, err := http.Post("http://localhost:7000/copy/second", m2.FormDataContentType(), pipeR2)

		if err != nil {
			log.Printf("Error on POST second: %v\n", err)
		}

		fmt.Printf("StatusCode second: %d\n", res.StatusCode)
	}()

	for {
		part, err := reader.NextPart()

		if err != nil {
			if err == io.EOF {
				fmt.Println("Main: finalizado as parts")
				break
			}

			log.Println(err)
			break
		}

		mw1, err := m1.CreateFormFile(part.FormName(), part.FileName())

		if err != nil {
			log.Printf("Error on CREATEFORM first: %v\n", err)
		}

		mw2, err := m2.CreateFormFile(part.FormName(), part.FileName())

		if err != nil {
			log.Printf("Error on CREATEFORM second: %v\n", err)
		}

		buf := make([]byte, 40)

		fmt.Printf("Main: lendo: %s\n", part.FormName())

		for {
			read, err := part.Read(buf)

			if err != nil && err != io.EOF {
				log.Println(err)
				break
			}

			if read == 0 {
				break
			}

			if _, err = mw1.Write(buf); err != nil {
				log.Printf("error on writing to first: %v", err)
			}

			if _, err = mw2.Write(buf); err != nil {
				log.Printf("error on writing to second: %v", err)
			}

			buf = make([]byte, 40)

			if err == io.EOF {
				break
			}
		}
	}

	fmt.Println("Closing m1 and m2")


	if err = m1.Close(); err != nil {
		log.Printf("Error on CLOSE first: %v\n", err)
	}

	if err = m2.Close(); err != nil {
		log.Printf("Error on CLOSE second: %v\n", err)
	}

	if err = pipeW1.Close(); err != nil {
		log.Printf("Error on CLOSEWPIPE first: %v\n", err)
	}

	if err = pipeW2.Close(); err != nil {
		log.Printf("Error on CLOSEWPIPE second: %v\n", err)
	}

	writeOk(w)
}

func handleFirst(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()

	if err != nil {
		log.Printf("First: is not multipart: %v\n", err)
		writeNotOk(w)

		return
	}

	for {
		part, err := reader.NextPart()

		if err != nil {
			if err == io.EOF {
				fmt.Println("First: finalizado as parts")
				break
			}

			log.Println(err)
			break
		}

		buf := make([]byte, 40)

		fmt.Printf("First: lendo: %s\n", part.FormName())

		for {
			read, err := part.Read(buf)

			if err != nil && err != io.EOF {
				log.Println(err)
				break
			}

			if read == 0 {
				break
			}

			fmt.Printf("First: %s\n", string(buf))

			buf = make([]byte, 40)

			if err == io.EOF {
				break
			}
		}
	}

	writeOk(w)
}

func handleSecond(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()

	if err != nil {
		log.Printf("Second: is not multipart: %v\n", err)
		writeNotOk(w)

		return
	}

	for {
		part, err := reader.NextPart()

		if err != nil {
			if err == io.EOF {
				fmt.Println("Second: finalizado as parts")
				break
			}

			log.Println(err)
			break
		}

		buf := make([]byte, 40)

		fmt.Printf("Second: lendo: %s\n", part.FormName())

		for {
			read, err := part.Read(buf)

			if err != nil && err != io.EOF {
				log.Println(err)
				break
			}

			if read == 0 {
				break
			}

			fmt.Printf("Second: %s\n", string(buf))

			buf = make([]byte, 40)

			if err == io.EOF {
				break
			}
		}
	}

	writeOk(w)
}

func writeOk(w http.ResponseWriter) {
	w.WriteHeader(200)

	if _, err := w.Write([]byte("ok")); err != nil {
		log.Printf("Error on write ok: %v", err)
	}
}

func writeNotOk(w http.ResponseWriter) {
	w.WriteHeader(500)

	if _, err := w.Write([]byte("not ok")); err != nil {
		log.Printf("Error on write nt ok: %v", err)
	}
}
