package main

import (
	"os"
	"flag"
	"strings"
	"fmt"
	"net/http"
	"math/rand"
	"sync"
)

var operators = make(map[string]bool)
var sessions = make(map[string]string)

var operatorsMutex = &sync.Mutex{}
var sessionsMutex = &sync.Mutex{}

func writeOperators (number string, state bool) {
	operatorsMutex.Lock()
	defer operatorsMutex.Unlock()
	operators[number] = state
}

func readOperators (number string) (bool, bool) {
	operatorsMutex.Lock()
	defer operatorsMutex.Unlock()
	op, exists := operators[number]
	return op, exists
}

// Write a nil-valued number to delete the session from the map
func writeSessions (sessionId string, number string) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()
	if number == "" {
		delete(sessions, sessionId)
	}else{
		sessions[sessionId] = number
	}
}

func readSessions (sessionId string) (string, bool) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()
	op, exists := sessions[sessionId]
	return op, exists
}

func main () {
	// Command-line flags
	var port, allOperators string
	flag.StringVar(&port, "port", "", "The port to bind to")
	flag.StringVar(&allOperators, "operators", "", "Comma-separated list of the operators to route calls to")
	flag.Parse()

	if port == "" || allOperators == "" {
		fmt.Println("Usage: ./app -port <port> -operators <Comma-separated list of operators>")
		os.Exit(0)
	}

	for _, op := range strings.Split(allOperators, ",") {
		operators[op] = false
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		active := r.FormValue("isActive")
		sessionId := r.FormValue("sessionId")
		op, sessionExists := readSessions(sessionId)
		
		if active == "0" {
			// Toggle operator's availability status
			writeOperators(op, false)
			writeSessions(sessionId, "")
		}else {
			// Create session
			if !sessionExists {
				inactiveOps := []string{}
				operatorsMutex.Lock()
				for op, isActive := range operators {
					if !isActive {
						inactiveOps = append(inactiveOps, op)
					}
				}
				operatorsMutex.Unlock()
				lenOps := len(inactiveOps)
				if lenOps == 0 {
					// All operators are busy
					fmt.Fprintf(w, `<Response><Say>Hello, all of our operators are currently busy, please call back in a bit.</Say></Response>`)
				}else{
					// There is at least one available operator
					var randomIndex int
					if lenOps > 1 {
						randomIndex = rand.Intn(lenOps - 1)
					}
					operator := inactiveOps[randomIndex]
					writeOperators(operator, true)
					writeSessions(sessionId, operator)
					fmt.Fprintf(w, `<Response><Dial phoneNumbers="%s"/></Response>`, operator)
				}
			}
		}
	})
	
	http.ListenAndServe(":"+port, nil)
}
