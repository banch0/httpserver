package main

import (
	"log"
	"net"
	"net/http"
	"testing"
	"time"
)

func start(host string) {
	listener, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Fatalf("Servet can't started error: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("can't accept client: ", err)
		}
		go requestHandler(conn)
	}
}

var netClient = &http.Client{
	Timeout: 50 * time.Second,
}

var link = "http://127.0.0.1:9999"

func Test_server(t *testing.T) {
	// naive approach
	go func() {
		listener, err := net.Listen("tcp", "0.0.0.0:9999")
		if err != nil {
			log.Fatalf("Servet can't started error: %v", err)
		}
		defer func() {
			err := listener.Close()
			if err != nil {
				t.Logf("Listener close error: %v", err)
			}
		}()

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("can't accept client: ", err)
			}
			go requestHandler(conn)
		}
	}()

	go func() {
		httpRequest, err := http.NewRequest("GET", link, nil)
		if err != nil {
			t.Error(err)
		}

		resp, err := netClient.Do(httpRequest)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Error(err)
		}
		time.Sleep(5 * time.Second)
	}()
}

func TestDownload(t *testing.T) {
	go func() {
		listener, err := net.Listen("tcp", "0.0.0.0:9999")
		if err != nil {
			log.Fatalf("Servet can't started error: %v", err)
		}
		defer func() {
			err := listener.Close()
			if err != nil {
				t.Logf("Listener close error: %v", err)
			}
		}()
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("can't accept client: ", err)
			}
			go requestHandler(conn)
		}
	}()

	go func() {
		httpRequest, err := http.NewRequest("GET", link, nil)
		if err != nil {
			t.Error(err)
		}

		resp, err := netClient.Do(httpRequest)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			t.Error(err)
		}
	}()
}

func TestImage(t *testing.T) {
	go func() {
		listener, err := net.Listen("tcp", "0.0.0.0:9999")
		if err != nil {
			log.Fatalf("Servet can't started error: %v", err)
		}
		defer func() {
			err := listener.Close()
			if err != nil {
				t.Logf("Listener close error: %v", err)
			}
		}()

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("can't accept client: ", err)
			}
			go requestHandler(conn)
		}
	}()

	go func() {
		httpRequest, err := http.NewRequest("GET", link+"/assets/image.jpg", nil)
		if err != nil {
			t.Error(err)
		}

		resp, err := netClient.Do(httpRequest)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Error(err)
		}
	}()
}

func TestImagePng(t *testing.T) {
	go func() {
		listener, err := net.Listen("tcp", "0.0.0.0:9999")
		if err != nil {
			log.Fatalf("Servet can't started error: %v", err)
		}
		defer func() {
			err := listener.Close()
			if err != nil {
				t.Logf("Listener close error: %v", err)
			}
		}()

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("can't accept client: ", err)
			}
			go requestHandler(conn)
		}
	}()

	go func() {
		httpRequest, err := http.NewRequest("GET", link+"/assets/image.png", nil)
		if err != nil {
			t.Error(err)
		}

		resp, err := netClient.Do(httpRequest)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Error(err)
		}
	}()
}

func TestPdf(t *testing.T) {
	go func() {
		listener, err := net.Listen("tcp", "0.0.0.0:9999")
		if err != nil {
			log.Fatalf("Servet can't started error: %v", err)
		}
		defer func() {
			err := listener.Close()
			if err != nil {
				t.Logf("Listener close error: %v", err)
			}
		}()

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("can't accept client: ", err)
			}
			go requestHandler(conn)
		}
	}()

	go func() {
		httpRequest, err := http.NewRequest("GET", link+"/assets/doc", nil)
		if err != nil {
			t.Error(err)
		}

		resp, err := netClient.Do(httpRequest)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Error(err)
		}
	}()
}

func TestTXT(t *testing.T) {
	go func() {
		listener, err := net.Listen("tcp", "0.0.0.0:9999")
		if err != nil {
			log.Fatalf("Servet can't started error: %v", err)
		}
		defer func() {
			err := listener.Close()
			if err != nil {
				t.Logf("Listener close error: %v", err)
			}
		}()

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("can't accept client: ", err)
			}
			go requestHandler(conn)
		}
	}()

	go func() {
		httpRequest, err := http.NewRequest("GET", link+"/assets/text", nil)
		if err != nil {
			t.Error(err)
		}

		resp, err := netClient.Do(httpRequest)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Error(err)
		}
	}()
}
