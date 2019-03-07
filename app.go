package main

import (
	"os"
	"flag"
	"strings"
	"fmt"
	"net/http"
	"math/rand"
)

var operators = make(map[string]bool)
var sessions = make(map[string]string)

func main () {
	// Command-line flags
	var port, allOperators, virtualNumber string
	flag.StringVar(&port, "port", "", "The port to bind to")
	flag.StringVar(&allOperators, "operators", "", "Comma-separated list of the operators to route calls to")
	flag.StringVar(&virtualNumber, "virtualNumber", "", "The virtualnumber on your AT account")
	flag.Parse()

	if port == "" || allOperators == "" {
		fmt.Println("Usage: ./app -port <port> -virtualNumber <Your virtualnumber> -operators <Comma-separated list of operators>")
		os.Exit(0)
	}

	for _, op := range strings.Split(allOperators, ",") {
		operators[op] = false
	}
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		active := r.FormValue("isActive")
		sessionId := r.FormValue("sessionId")
		callerNumber := r.FormValue("callerNumber")
		op, sessionExists := sessions[sessionId]
		
		// Check if it's an operator calling and dequeue
		// If the queue is empty this will just hang up, which is fine
		if _,exists := operators[callerNumber]; exists {
			fmt.Fprintf(w, `<Response><Dequeue phoneNumber="%s"/></Response>`, virtualNumber)
			return
		}

		if active == "0" {
			// Toggle operator's availability status
			if sessionExists {
				operators[op] = false
				delete(sessions, sessionId)
			}
		}else {
			// Create session
			if !sessionExists {
				inactiveOps := []string{}
				for op, isActive := range operators {
					if !isActive {
						inactiveOps = append(inactiveOps, op)
					}
				}
				lenOps := len(inactiveOps)
				if lenOps == 0 {
					// All operators are busy
					fmt.Fprintf(w, `<Response><Enqueue/></Response>`)
				}else{
					// There is at least one available operator
					var randomIndex int
					if lenOps > 1 {
						randomIndex = rand.Intn(lenOps - 1)
					}
					operator := inactiveOps[randomIndex]
					operators[operator] = true
					sessions[sessionId] = operator
					fmt.Fprintf(w, `<Response><Dial phoneNumbers="%s"/></Response>`, operator)
				}
			}
		}
	})
	
	http.ListenAndServe(":"+port, nil)
}
