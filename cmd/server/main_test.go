package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleFile(t *testing.T) {
	dir := http.Dir("./assets/")

	// strip prefix instead /temp/*.jpg /*.jpg
	handler := http.StripPrefix("/assets", http.FileServer(dir))
	mux := http.NewServeMux()
	mux.Handle("/assets/", handler)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/assets/", nil)
	mux.ServeHTTP(writer, request)
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
}
