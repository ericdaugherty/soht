package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net"
	"net/http"
	"strconv"
)

// Statu
const StatusBadParam = 400

var homeTempl = template.Must(template.New("home").Parse(homeStr))
var adminTempl = template.Must(template.New("home").Parse(adminStr))

var connections = make(map[string] net.Conn)

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
	conn.Close()
	return OpenResponse{ true, connId, "" }
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
