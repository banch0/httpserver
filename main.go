package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:9898")
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

	var parts []string
	headers := make(map[string]string)
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

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

		switch {
		case method == "GET" && request == "/" && protocol == "HTTP/1.1":
			data, err := ioutil.ReadFile("./assets/main.html")
			if err != nil {
				log.Println("Error read file", err)
			}
			err = sendResponse(response, writer, data, "text/html")
			if err != nil {
				log.Println(err)
			}
		case method == "GET" && request == "/assets/doc":
			data, err := ioutil.ReadFile("./assets/doc.pdf")
			if err != nil {
				log.Println("Error read file", err)
			}
			err = sendResponse(response, writer, data, "application/pdf")
			if err != nil {
				log.Println(err)
			}
		case method == "GET" && request == "/assets/doc.pdf?download":
			data, err := ioutil.ReadFile("./assets/doc.pdf")
			if err != nil {
				log.Println("Error read file", err)
			}

			err = downloadFile(response, writer, data, "application/pdf")
			if err != nil {
				log.Println(err)
			}
		case method == "GET" && request == "/assets/image.jpg":
			data, err := ioutil.ReadFile("./assets/image.jpg")
			if err != nil {
				log.Println("Error read file", err)
			}

			err = sendResponse(response, writer, data, "image/jpg")
			if err != nil {
				log.Println(err)
			}
		case method == "GET" && request == "/assets/image.jpg?download":
			data, err := ioutil.ReadFile("./assets/image.jpg")
			if err != nil {
				log.Println("Error read file", err)
			}

			err = downloadFile(response, writer, data, "image/jpg")
			if err != nil {
				log.Println(err)
			}
		case method == "GET" && request == "/assets/music":
			data, err := ioutil.ReadFile("./assets/miyaGi.mp3")
			if err != nil {
				log.Println("Error read file", err)
			}
			err = sendResponse(response, writer, data, "audio/mpeg")
			if err != nil {
				log.Println(err)
			}
		case method == "GET" && request == "/assets/miyaGi.mp3?download":
			data, err := ioutil.ReadFile("./assets/miyaGi.mp3")
			if err != nil {
				log.Println("Error read file", err)
			}

			err = downloadFile(response, writer, data, "audio/mpeg")
			if err != nil {
				log.Println(err)
			}
		case method == "GET" && request == "/assets/video":
			data, err := ioutil.ReadFile("./assets/index.webm")
			if err != nil {
				log.Println("Error read file", err)
			}
			err = sendResponse(response, writer, data, "video/webm")
			if err != nil {
				log.Println(err)
			}
		case method == "GET" && request == "/assets/index.webm?download":
			data, err := ioutil.ReadFile("./assets/index.webm")
			if err != nil {
				log.Println("Error read file", err)
			}

			err = downloadFile(response, writer, data, "video/webm")
			if err != nil {
				log.Println(err)
			}
		case method == "GET" && request == "/assets/image.png":
			data, err := ioutil.ReadFile("./assets/biohazard.png")
			if err != nil {
				log.Println("ReadFile Error: ", err)
			}
			err = sendResponse(response, writer, data, "image/jpg")
			if err != nil {
				log.Println(err)
			}
		case method == "GET" && request == "/assets/image.png?download":
			data, err := ioutil.ReadFile("./assets/biohazard.png")
			if err != nil {
				log.Println("Error read file", err)
			}

			err = downloadFile(response, writer, data, "image/png")
			if err != nil {
				log.Println(err)
			}
		case method == "GET" && request == "/assets/text":
			data, err := ioutil.ReadFile("./assets/myfile.txt")
			if err != nil {
				log.Println("Error read file", err)
			}
			err = sendResponse(response, writer, data, "text/plain")
			if err != nil {
				log.Println(err)
			}
		case method == "GET" && request == "/assets/myfile.txt?download":
			data, err := ioutil.ReadFile("./assets/myfile.txt")
			if err != nil {
				log.Println("Error read file", err)
			}

			err = downloadFile(response, writer, data, "text/plain")
			if err != nil {
				log.Println(err)
			}
		case method == "GET" && request == "/favicon.ico":
			writer.WriteString("HTTP/1.1 200 OK\r\n")
			writer.WriteString("\r\n")
			writer.Flush()
			log.Println("favicon")
		default:
			data, err := ioutil.ReadFile("./assets/404.html")
			if err != nil {
				log.Println("Error read file", err)
			}

			err = sendResponse(response, writer, data, "text/html")
			if err != nil {
				log.Println(err)
			}
		}
	}

	for key, value := range headers {
		log.Printf("key: %s value: %s", key, value)
	}
}

func sendResponse(resp bytes.Buffer, writer *bufio.Writer, data []byte, content string) error {
	resp.WriteString("HTTP/1.1 200 OK\r\n")
	resp.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(data)))
	resp.WriteString("Connection: Close\r\n")
	resp.WriteString("Content-Type: " + content + "\r\n")
	resp.WriteString("\r\n")
	resp.Write(data)
	resp.WriteTo(writer)
	err := writer.Flush()
	if err != nil {
		log.Printf("Writing response error: %v", err)
	}
	return err
}

func downloadFile(resp bytes.Buffer, writer *bufio.Writer, data []byte, content string) error {
	resp.WriteString("HTTP/1.1 200 OK\r\n")
	resp.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(data)))
	resp.WriteString("Content-Type: application/octet-stream\r\n")
	resp.WriteString("Content-Transfer-Encoding: binary\r\n")
	resp.WriteString("\r\n")
	resp.Write(data)
	resp.WriteTo(writer)
	err := writer.Flush()
	if err != nil {
		log.Printf("Writing response error: %v", err)
	}
	return err
}
