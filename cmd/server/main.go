package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"time"
)

var port = ":9000"

func main() {

	// name of serving dir
	dir := http.Dir("./assets/")

	// strip prefix instead /temp/*.jpg /*.jpg
	handler := http.StripPrefix("/assets", http.FileServer(dir))
	http.Handle("/assets/", handler)

	// api for downloading one file
	http.HandleFunc("/original/", downloadFile)

	// starting server
	log.Println("server starting on localhost" + port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalln(err)
	}
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	filePath := path.Base(r.URL.Path)
	data, err := ioutil.ReadFile("./assets/" + filePath)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	http.ServeContent(w, r, "./assets/"+filePath, time.Now(), bytes.NewReader(data))
}
