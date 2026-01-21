package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type logEnt struct {
	timeStamp   string
	randStr     string
	message     string
	fileContent string
}

var (
	message     string
	fileContent string
	pingPort    string
)

func init() {
	pingPort = os.Getenv("PING_PORT")
	if pingPort == "" {
		fmt.Println("env PING_PORT was unset\nUsing Port 3001 as pingPort")
		pingPort = "3001"
	}

	message = os.Getenv("MESSAGE")
	if message == "" {
		log.Print("env MESSAGE was unset")
	}

	f := os.Getenv("INFO_FILE")
	if f == "" {
		log.Print("env INFO_FILE was unset")
	} else {
		data, err := os.ReadFile(f)
		if err != nil {
			log.Printf("Read INFO_FILE Error: %v", err)
		}
		fileContent = string(data)
	}
}

func (l logEnt) String() string {
	return fmt.Sprintf("%s\n%s\n%s %s\n",
		l.message, l.fileContent, l.timeStamp, l.randStr)
}

func randomString() string {
	b := make([]byte, 4) // 4 bytes = 8 hex characters
	if _, err := rand.Read(b); err != nil {
		return "00000000"
	}
	return hex.EncodeToString(b)
}

func main() {
	logPort := os.Getenv("PORT")
	if logPort == "" {
		fmt.Println("env PORT was unset\nUsing Port 3000 as logPort")
		logPort = "3000"
	}
	addr := ":" + logPort
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)

	log.Printf("Server starting on %s\n", addr)
	log.Fatal((http.ListenAndServe(addr, mux)))
}

var myClient = &http.Client{
	Timeout: time.Second * 5,
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	hostAddr := "http://localhost:" + pingPort + "/pings"

	resp, err := myClient.Get(hostAddr)
	if err != nil {
		//  502 for upstream errors
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("bad status: %s", resp.Status), resp.StatusCode)
		return
	}
	lr := io.LimitReader(resp.Body, 1024*1024)
	data, err := io.ReadAll(lr)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}
	var count int
	_, err = fmt.Sscanf(string(data), "ping %d", &count)
	if err != nil {
		count = 0
	}
	ent := logEnt{
		timeStamp:   time.Now().Format("2006-01-02 15:04:05"),
		randStr:     randomString(),
		message:     message,
		fileContent: fileContent,
	}
	fmt.Fprintf(w, "%v Ping / Pongs: %d", ent, count)
}
