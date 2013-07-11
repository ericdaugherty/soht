package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
)

var initialized bool

func initServer() {
	if !initialized {
		initialized = true
		go runServer(":8100")
	}
}

func TestOpenNoHost(t *testing.T) {

	initServer()
	urlVals := make(url.Values)
	urlVals.Set("action", "open")
	urlVals.Set("body", `{"Port":80}`)
	resp, err := http.PostForm("http://localhost:8100/", urlVals)

	if err != nil {
		t.Errorf("Unexpected error", err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("Recieved response code %v but should be 400.  Host was not sent", resp.StatusCode)
	}
}

func TestOpenNoPort(t *testing.T) {

	initServer()
	urlVals := make(url.Values)
	urlVals.Set("action", "open")
	urlVals.Set("body", `{"Host":"google.com"}`)

	resp, err := http.PostForm("http://localhost:8100/", urlVals)

	if err != nil {
		t.Errorf("Unexpected error", err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("Recieved response code %v but should be 400.  Port was not sent", resp.StatusCode)
	}
}

func TestOpenInvalidPort(t *testing.T) {

	initServer()
	urlVals := make(url.Values)
	urlVals.Set("action", "open")
	urlVals.Set("body", `{"Host":"google.com","Port":"abcd"}`)

	resp, err := http.PostForm("http://localhost:8100/", urlVals)

	if err != nil {
		t.Errorf("Unexpected error", err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("Recieved response code %v but should be 400.  Host was not sent", resp.StatusCode)
	}
}

func TestOpenSimple(t *testing.T) {

	initServer()
	urlVals := make(url.Values)
	urlVals.Set("action", "open")
	urlVals.Set("body", `{"Host":"google.com","Port":80}`)

	resp, err := http.PostForm("http://localhost:8100/", urlVals)
	respBody, _ := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("Error opening connection to Google")
	} else if resp.StatusCode != 200 {
		t.Errorf("Error opening connection, invalid status code.  Body:", string(respBody))
	} else {
		var openResponse OpenResponse
		_ = json.Unmarshal(respBody, &openResponse)
		fmt.Println("Response", openResponse)
		if !openResponse.Success {
			t.Errorf("Server unable to open connection. Desc: %v", openResponse.ErrorDesc)
		}
	}
}

func BenchmarkOpenConnectionGoogle(b *testing.B) {
	noopWriter := NoopWriter{}	
	log.SetOutput(noopWriter)

	for n := 0; n < b.N; n++ {
		openConnection("google.com", 80)
	}
	
	log.SetOutput(os.Stderr)
}
