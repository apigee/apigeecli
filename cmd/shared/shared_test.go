package shared

import (
	"testing"
)

func TestInit(t *testing.T) {
	Init()
	Info.Println("Printing Information")
	Warning.Println("Printing Warning")
	Error.Println("Printing Error")
}

func TestHttpGet(t *testing.T) {
	Init()
	err := HttpClient("https://httpbin.org/get")
	if err != nil {
		t.Fatal(err)
	}
}

func TestHttpPost(t *testing.T) {
	payload := "test"
	err := HttpClient("https://httpbin.org/post", payload)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHttpDelete(t *testing.T) {
	err := HttpClient("https://httpbin.org/post", "", "DELETE")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDownloadResource(t *testing.T) {
	//download 1000 bytes
	Init()
	err := DownloadResource("https://httpbin.org/stream-bytes/1000", "test")
	if err != nil {
		t.Fatal(err)
	}
}
