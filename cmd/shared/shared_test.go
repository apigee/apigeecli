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
	_, err := HttpClient(true, "https://httpbin.org/get")
	if err != nil {
		t.Fatal(err)
	}
}

func TestHttpPost(t *testing.T) {
	payload := "test"
	_, err := HttpClient(true, "https://httpbin.org/post", payload)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHttpDelete(t *testing.T) {
	_, err := HttpClient(true, "https://httpbin.org/delete", "", "DELETE")
	if err != nil {
		t.Fatal(err)
	}
}

func TestHttpInvalidParam(t *testing.T) {
	_, err := HttpClient(true, "https://httpbin.org/delete", "", "DELTE")
	if err == nil {
		t.Fatal(err)
	}
}

func TestHttpInvalidNumberOfParams(t *testing.T) {
	_, err := HttpClient(true, "https://httpbin.org/delete", "", "DELETE", "SOMETHING ELSE")
	if err == nil {
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

func TestWriteJSONArrayToFile(t *testing.T) {
	Init()
	var entityPayloadList = []byte{'g', 'o', 'l', 'a', 'n', 'g'}
	err := WriteByteArrayToFile("test.json", false, entityPayloadList)
	if err != nil {
		t.Fatal(err)
	}
}
