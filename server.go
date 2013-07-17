package main

import (
	"bufio"
	"encoding/json"
	"html/template"
	"log"
	"net"
	"net/http"
	"strconv"
)

const StatusBadParam = http.StatusBadRequest
const StatusNotFound = http.StatusNotFound

var homeTempl = template.Must(template.New("home").Parse(homeStr))
var adminTempl = template.Must(template.New("home").Parse(adminStr))

var connections = make(map[uint32] ConnectionInfo)

var counter = make(chan uint32, 1)

type OpenRequest struct {
	Host string
	Port uint16
	Username string
	Password string
}

type OpenResponse struct {
	Success bool
	ConnectionId uint32
	ErrorDesc string
}

type ReadRequest struct {
	ConnectionId uint32
}

type ConnectionInfo struct {
	ConnectionId uint32
	Connection net.Conn
}

func init() {
	go connectionIdCounter(counter)
}

func runServer(serveraddr string) {
	http.Handle("/", http.HandlerFunc(rootHandler))
	http.Handle("/admin/", http.HandlerFunc(adminHandler))
	err := http.ListenAndServe(serveraddr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" && req.URL.Path == "/" {
		http.Redirect(w, req, "/admin/", 302)
		return
	}
	if req.Method == "GET" {
		log.Printf("Request made to %v returning 404.\n", req.URL)
		http.NotFound(w, req)
		return
	}
	req.ParseForm()
	log.Printf("Request Parms: %v\n", req.Form)
	values := req.PostForm
	action := values.Get("action")
	jsonBody := values.Get("body")
	switch action {
	case "open":
		var openRequest OpenRequest 
		jsonError := json.Unmarshal([]byte(jsonBody), &openRequest)	
		if jsonError != nil {
			http.Error(w, "Invalid JSON in body parameter. " + jsonError.Error(), StatusBadParam)
			return
		}
		if len(openRequest.Host) == 0 {
			http.Error(w, "host parameter not specified", StatusBadParam)
			return
		}
		if  openRequest.Port == 0 {
			http.Error(w, "port parameter not specified", StatusBadParam)
			return
		}

		openResponse := openConnection(openRequest.Host, openRequest.Port)
		responseBody, jsonError := json.Marshal(openResponse)

		// Backup Response in case the Marshal call fails
		if jsonError != nil {
			responseBody = []byte(`{"Success":false,"ConnectionId":0,"ErrorDesc":"Unable to parse response."}`)
		}
		
		w.Write(responseBody)
	case "read":
		var readRequest ReadRequest 
		jsonError := json.Unmarshal([]byte(jsonBody), &readRequest)	
		if jsonError != nil {
			http.Error(w, "Invalid JSON in body parameter. " + jsonError.Error(), StatusBadParam)
			return
		}
		connInfo, ok := connections[readRequest.ConnectionId]
		if !ok {
			log.Println("Request for unknown connectionId:", readRequest.ConnectionId)
			http.Error(w, "Unknown Connection Id.", StatusNotFound)	
			return
		}
		longRead(w, connInfo)
	default:
		log.Println("Unknown Command")
		http.Error(w, "Unknown Command", 400)
	}
}

func adminHandler(w http.ResponseWriter, req *http.Request) {
	adminTempl.Execute(w, "")
}

func openConnection(host string, port uint16) OpenResponse {

	var address string = host + ":" + strconv.FormatUint(uint64(port), 10)
	log.Println("Opening address:" , address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Println("Open Failed", err)
		return OpenResponse { false, 0, err.Error() }
	}
	connId := <-counter
	connInfo := ConnectionInfo { connId, conn }
	connections[connId] = connInfo;
	return OpenResponse{ true, connId, "" }
}

func longRead(w http.ResponseWriter, connInfo ConnectionInfo) {
	log.Println("Starting Long Read for connectionId:", connInfo.ConnectionId)
	conn := connInfo.Connection
	bytesIn := make([]byte, 512)
	bufReader := bufio.NewReader(conn)

	w.WriteHeader(http.StatusOK)
	for {
		log.Println("Looping in longRead")
		bytesRead, error := bufReader.Read(bytesIn)
		log.Printf("Read %v bytes\r\n", bytesRead)
		if error != nil {
			log.Println("Error reading from connection, ending read goroutine.")
			return
		}
		w.Write(bytesIn)
		return;
	}
}

func connectionIdCounter(c chan<- uint32) {
        var counter uint32 = 0
        for {
                counter++
		// Reset before we overflow, assuming 32 bit.
		if counter > 4294967200 {
			counter = 1
		}
                c <- counter 
        }
}

const homeStr = `<html><body>Hello from SOHT</body></html>`
const adminStr = `<html><body>Hello from SOHT Admin</body></html>`
