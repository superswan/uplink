package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
)

// Client struct to keep track of each Client's shell session
type Client struct {
	conn net.Conn `json:"-"`
	Ip   string   `json:"ip"`
	Port string   `json:"port"`
	Id   string   `json:"id"`
	DbId string   `json:"dbId"`
}

// record struct for interacting with DB
type Record struct {
	ClientID     string `json:"client_id"`
	IPAddress    string `json:"ipaddress"`
	Port         string `json:"port"`
	SessionCount int    `json:"sessionCount"`
}

type Response struct {
	Page       int    `json:"page"`
	PerPage    int    `json:"perPage"`
	TotalPages int    `json:"totalPages"`
	TotalItems int    `json:"totalItems"`
	Items      []Item `json:"items"`
}

type Item struct {
	Id             string `json:"id"`
	CollectionId   string `json:"collectionId"`
	CollectionName string `json:"collectionName"`
	Created        string `json:"created"`
	Updated        string `json:"updated"`
	ClientID       string `json:"client_id"`
	IPAddress      string `json:"ipaddress"`
	Port           string `json:"port"`
	SessionCount   int    `json:"sessionCount"`
}

type Command struct {
	Command string `json:"command"`
}

type Output struct {
	Output string `json:"output"`
}

var clients []*Client
var activeClient *Client
var mutex = &sync.Mutex{}

// Routes output to different channels based on CLI or API requests
var messages = make(map[string]chan string)

func upsertClient(c *Client) {
	url := fmt.Sprintf("http://127.0.0.1:8090/api/collections/clients/records?filter=(client_id='%s')", c.Id)

	// Check if client exists
	getReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	getResp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer getResp.Body.Close()

	var response Response
	json.NewDecoder(getResp.Body).Decode(&response)

	if len(response.Items) == 0 {
		// Client does not exist, create a new record
		record := &Record{
			ClientID:     c.Id,
			IPAddress:    c.Ip,
			Port:         c.Port,
			SessionCount: 1,
		}

		jsonData, err := json.Marshal(record)
		if err != nil {
			fmt.Println(err)
			return
		}

		postReq, err := http.NewRequest("POST", "http://127.0.0.1:8090/api/collections/clients/records", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println(err)
			return
		}

		postReq.Header.Set("Content-Type", "application/json")

		postResp, err := http.DefaultClient.Do(postReq)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer postResp.Body.Close()

		body, err := ioutil.ReadAll(postResp.Body)
		if err != nil {
			fmt.Println("Error reading body:", err)
			return
		}

		var postResponse Item
		json.Unmarshal(body, &postResponse)

		c.DbId = postResponse.Id

	} else {
		// Client exists, update the record
		existingRecord := response.Items[0]

		// Create a new struct for updating the record
		type UpdateRecord struct {
			IPAddress    string `json:"ipaddress"`
			Port         string `json:"port"`
			SessionCount int    `json:"sessionCount"`
		}

		updateRecord := &UpdateRecord{
			IPAddress:    c.Ip,
			Port:         c.Port,
			SessionCount: existingRecord.SessionCount + 1,
		}

		jsonData, err := json.Marshal(updateRecord)
		if err != nil {
			fmt.Println(err)
			return
		}

		patchReq, err := http.NewRequest("PATCH", "http://127.0.0.1:8090/api/collections/clients/records/"+existingRecord.Id, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println(err)
			return
		}

		patchReq.Header.Set("Content-Type", "application/json")

		_, err = http.DefaultClient.Do(patchReq)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func handleConnection(c *Client) {
	reader := bufio.NewReader(c.conn)
	msgChan := make(chan string, 1)
	messages[c.Id] = msgChan

	for {
		lenBuf := make([]byte, 4)
		_, err := io.ReadFull(reader, lenBuf)
		if err != nil {
			fmt.Printf("Error reading from client %s: %v\n", c.Id, err)
			close(msgChan)
			delete(messages, c.Id)
			break
		}

		msgLen := binary.BigEndian.Uint32(lenBuf)
		msgBuf := make([]byte, msgLen)

		_, err = io.ReadFull(reader, msgBuf)
		if err != nil {
			fmt.Printf("Error reading from client %s: %v\n", c.Id, err)
			close(msgChan)
			delete(messages, c.Id)
			break
		}

		message := string(msgBuf)
		msgChan <- message

		if message == "exit" {
			break
		}

		mutex.Lock()
		if activeClient == c {
			output := <-msgChan
			fmt.Printf("%s: %s\n", c.Id, output)
			fmt.Print("> ")
		}
		mutex.Unlock()
	}

	fmt.Printf("Client %s disconnected\n", c.Id)
	c.conn.Close()

	for i := range clients {
		if clients[i] == c {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

func enableCors(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fn(w, r)
	}
}

func handleCommand(w http.ResponseWriter, r *http.Request) {
	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusBadRequest)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id := parts[2]
	if len(id) == 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var cmd Command
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	// Find client by ID
	var matchingClient *Client
	for _, c := range clients {
		if c.Id == id {
			matchingClient = c
			break
		}
	}

	if matchingClient == nil {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	// Send command to the shell here using id and get output.
	if matchingClient != nil {
		_, err := matchingClient.conn.Write([]byte(cmd.Command + "\n"))
		if err != nil {
			fmt.Printf("Error sending command to client %s: %v\n", matchingClient.Id, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Receive output from shell.
		output := <-messages[matchingClient.Id]

		// Send the output back in response.
		response := Output{
			Output: output,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}

	} else {
		fmt.Println("No matching client found.")
		w.WriteHeader(http.StatusNotFound)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// Lock the mutex before accessing the shared clients slice
	mutex.Lock()
	defer mutex.Unlock()
	// Filter the slice to only include active clients' IDs.
	activeClientIDs := []string{}
	for _, client := range clients {
		activeClientIDs = append(activeClientIDs, client.Id)
	}

	// Now we marshal only the active client IDs to JSON.
	response, _ := json.Marshal(activeClientIDs)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
	} else {
		w.Write(response)
	}
}

func startAPI() {
	http.HandleFunc("/command/", enableCors(handleCommand))
	http.HandleFunc("/status/", statusHandler)

	fmt.Println("API server listening on localhost:8088")
	if err := http.ListenAndServe(":8088", nil); err != nil {
		log.Fatal(err)
	}
}

func main() {
	go startAPI()
	l, err := net.Listen("tcp", "0.0.0.0:5550")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer l.Close()
	fmt.Println("Server started...")
	fmt.Println("\nEnter 'show' to list clients, a client id to switch to that client, or a command to send to the active client:")

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				fmt.Print("> ")
				continue
			}

			// Split the RemoteAddr into IP and port
			ip, port, _ := net.SplitHostPort(conn.RemoteAddr().String())

			// Read client id
			clientID, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println("Error reading client id: ", err.Error())
				fmt.Print("> ")
				continue
			}
			clientID = strings.TrimSpace(clientID)

			c := &Client{
				conn: conn,
				Ip:   ip,
				Port: port,
				Id:   clientID,
			}

			// Upsert Client
			upsertClient(c)

			clients = append(clients, c)

			fmt.Printf("Client %s connected\n", c.Id)
			fmt.Print("> ")

			go handleConnection(c)
		}
	}()

	// CLI
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "show" {
			for _, c := range clients {
				fmt.Printf("Client ID: %s, IP: %s, Port: %s\n", c.Id, c.Ip, c.Port)
			}
		} else {
			// Check if the input matches a client id
			var matchingClient *Client
			for _, c := range clients {
				if c.Id == input {
					matchingClient = c
					break
				}
			}

			if matchingClient != nil {
				// Switch to the specified client
				mutex.Lock()
				activeClient = matchingClient
				mutex.Unlock()
				fmt.Printf("Switched to client %s\n", activeClient.Id)
			} else {
				// Send the input as a command to the active client
				if activeClient != nil {
					_, err := activeClient.conn.Write([]byte(input + "\n"))
					if err != nil {
						fmt.Printf("Error sending command to client %s: %v\n", activeClient.Id, err)
					}
				} else {
					fmt.Println("No active client. Please switch to a client first.")
				}
			}
		}
	}
}
