package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	urlVals.Set("body", `{"Host":"gmail-smtp-in.l.google.com","Port":25}`)

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

func TestReadUnknownConnection(t *testing.T) {

	initServer()
	urlVals := make(url.Values)
	urlVals.Set("action", "read")
	urlVals.Set("body", `{"ConnectionId":0}`)

	resp, err := http.PostForm("http://localhost:8100/", urlVals)
	respBody, _ := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("Error reading connection")
	} else if resp.StatusCode != 404 {
		t.Errorf("Reading of unknown connection should fail with 404.  Body:", string(respBody))
	} 
}


func TestReadOpenedConnection(t *testing.T) {

	initServer()

	initServer()
	urlVals := make(url.Values)
	urlVals.Set("action", "read")
	urlVals.Set("body", `{"ConnectionId":1}`)
	req, err := http.NewRequest("POST", "http://localhost:8100/", strings.NewReader(urlVals.Encode()))
	if err != nil { 
		t.Errorf("Error building Post")
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	conn, err := net.Dial("tcp", "localhost:8100")
	err = req.Write(conn)
	if err != nil {
		t.Errorf("Error writring request")
		return
	}

	bufReader := bufio.NewReader(conn)
	bytesIn := make([]byte, 512)
	var resString string
	for {
		bytesRead, error := bufReader.Read(bytesIn)
		fmt.Println("", bytesIn)
		if error != nil {
			t.Errorf("Error reading bytes from server response. %v", error)
			return
		}
		resString = string(bytesIn[0:bytesRead])

		_, error = http.ReadResponse(bufio.NewReader(strings.NewReader(resString)), req)

		if error != nil {
			t.Errorf("Error creating Response struct %v", error)
			return
		}

		conn.Close()
		break
	}

	respBody := strings.Split(resString, "\r\n\r\n")[1]
	if index := strings.Index(respBody, "220"); index != 0 {
		t.Errorf("Index of %v not expected.  Should be 0", index)
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
