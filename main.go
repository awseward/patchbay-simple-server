package main

import (
	"io"
	"log"
	"net/http"
	"sync"
)

type Stream struct {
	reader io.ReadCloser
	done   chan struct{}
}

func main() {
	var portStr string
	portStr = "9001"

	log.Println("Starting up: localhost" + ":" + portStr)

	channels := make(map[string]chan Stream)
	mutex := &sync.Mutex{}

	handler := func(w http.ResponseWriter, r *http.Request) {

		mutex.Lock()

		_, ok := channels[r.URL.Path]
		if !ok {
			channels[r.URL.Path] = make(chan Stream)
		}

		channel := channels[r.URL.Path]

		mutex.Unlock()

		log.Println(channel)
		if r.Method == "GET" {

			select {
			case stream := <-channel:
				io.Copy(w, stream.reader)
				close(stream.done)
			case <-r.Context().Done():
				log.Println("consumer canceled")
			}
		} else if r.Method == "POST" {

			doneSignal := make(chan struct{})
			stream := Stream{reader: r.Body, done: doneSignal}
			select {
			case channel <- stream:
				log.Println("connected to consumer")
			case <-r.Context().Done():
				log.Println("producer canceled")
			}

			<-doneSignal
		}
	}

	err := http.ListenAndServe(":"+portStr, http.HandlerFunc(handler))
	if err != nil {
		log.Fatal(err)
	}
}
