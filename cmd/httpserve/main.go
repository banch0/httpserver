package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

var assets = "./assets/"
var address = "0.0.0.0"

var extensions = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".html": "text/html; charset=utf-8",
	".txt":  "text/plain",
	".css":  "text/css; charset=utf-8",
	".pdf":  "application/pdf",
	".png":  "image/png",
	".webm": "video/webm",
	".mp3":  "audio/mpeg",
	".ico":  "image/x-icon",
}

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "9898"
	}
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		log.Fatalf("Servet can't started error: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("can't accept client: ", err)
		}
		log.Println("client connected")
		go requestHandler(conn)
	}
}

func requestHandler(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	headers := make(map[string]string)
	parts := make([]string, 0)

	for {
		readLine, err := reader.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("requestHandler ReadSlice Error: ", err)
		}

		if string(readLine) == "\r\n" {
			break
		}

		splited := strings.Split(string(readLine), ": ")

		if len(splited) > 1 {
			headers[splited[0]] = splited[1]
		} else {
			parts = strings.Split(strings.TrimSpace(string(readLine)), " ")
			if len(parts) != 3 {
				return
			}
		}
	}

	if len(parts) > 2 {
		method, request, protocol := parts[0], parts[1], parts[2]
		var response = bytes.Buffer{}

		if protocol != "HTTP/1.1" {
			return // errors.New("Wrong protocol")
		}

		if method != "GET" {
			return // errors.New("Method Not Allowed")
		}

		switch {
		case strings.HasPrefix(request, "/assets/"):

			filename, extension, download := handler(request)

			data, err := ioutil.ReadFile(assets + filename)
			if err != nil {
				log.Println("Error read file", err)
			}

			err = Response(response, writer, data, download, extension)
			if err != nil {
				log.Println("Send Response Error: ", err)
			}

		case request == "/":
			data, err := ioutil.ReadFile(assets + "main.html")
			if err != nil {
				log.Println("Error read file", err)
			}

			err = Response(response, writer, data, false, "text/html")
			if err != nil {
				log.Println(err)
			}

		case request == "/favicon.ico":
			writer.WriteString("HTTP/1.1 200 OK\r\n")
			writer.WriteString("\r\n")
			writer.Flush()

		default:
			data, err := ioutil.ReadFile(assets + "404.htm")
			if err != nil {
				log.Println("Error read file", err)
			}

			err = Response(response, writer, data, false, "text/html")
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// Response ...
func Response(resp bytes.Buffer,
	writer *bufio.Writer,
	data []byte,
	flag bool,
	contentType string) error {
	resp.WriteString("HTTP/1.1 200 OK\r\n")
	resp.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(data)))

	if flag {
		resp.WriteString("Content-Type: application/octet-stream\r\n")
		resp.WriteString("Content-Transfer-Encoding: binary\r\n")
	} else {
		resp.WriteString("Connection: Close\r\n")
		resp.WriteString("Content-Type: " + contentType + "\r\n")
	}
	
	resp.WriteString("\r\n")
	resp.Write(data)
	resp.WriteTo(writer)
	err := writer.Flush()
	if err != nil {
		log.Printf("Writing response error: %v", err)
	}
	return err

}

func handler(request string) (string, string, bool) {
	var download bool

	if strings.Contains(request, "?download") {
		download = true
		request = strings.TrimSuffix(request, "?download")
	}

	fileName := strings.TrimPrefix(request, "/assets/")
	extension := strings.Split(fileName, ".")

	return fileName, extensions["."+extension[1]], download
}
