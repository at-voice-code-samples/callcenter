package main

import (
	"os"
	"flag"
	"strings"
	"fmt"
	"net/http"
	"math/rand"
	"sync"
	"net/url"
	"bytes"
	"github.com/chzyer/readline"
)

var apiKey string
var userName string
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

func transfer (sessionId string, operatorNumber string) {
	operatorsMutex.Lock()
	if operators[operatorNumber] {
		fmt.Printf("Could not transfer call to %s: operator is on call.\n", operatorNumber)
		return
	}
	operatorsMutex.Unlock()
	
	reqBody := fmt.Sprintf("username=%s&phoneNumber=%s&callLeg=callee&sessionId=%s", url.QueryEscape(userName), url.QueryEscape(operatorNumber), url.QueryEscape(sessionId))
	req, err := http.NewRequest(http.MethodPost, "https://voice.africastalking.com/callTransfer", strings.NewReader(reqBody))
	if err != nil {
		fmt.Println("Could not create a new request:", err)
		return
	}
	// req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Apikey", apiKey)
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("Call transfer HTTP request failed:", err)
		return
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	fmt.Printf("Response: %s\n", buf.String())
}

func main () {
	// Command-line flags
	var port, allOperators, virtualNumber, apiKeyArg, userNameArg string
	flag.StringVar(&port, "port", "", "The port to bind to")
	flag.StringVar(&allOperators, "operators", "", "Comma-separated list of the operators to route calls to")
	flag.StringVar(&virtualNumber, "virtualNumber", "", "The virtualnumber on your AT account")
	flag.StringVar(&apiKeyArg, "apikey", "", "Your AT Apikey")
	flag.StringVar(&userNameArg, "username", "", "Your AT username")
	flag.Parse()

	apiKey = apiKeyArg
	userName = userNameArg

	if port == "" || allOperators == "" || virtualNumber == "" || apiKey == "" || userName == ""{
		fmt.Println("Usage: ./app -port <port> -virtualNumber <Your virtualnumber> -username <Your AT account username>",
			    "-apikey <Your API key> -operators <Comma-separated list of operators>")
		os.Exit(0)
	}

	for _, op := range strings.Split(allOperators, ",") {
		operators[op] = false
	}
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		active := r.FormValue("isActive")
		sessionId := r.FormValue("sessionId")
		callerNumber := r.FormValue("callerNumber")
		op, sessionExists := readSessions(sessionId)
		// Check if it's an operator calling and dequeue
		// If the queue is empty this will just hang up, which is fine
		if _,exists := readOperators(callerNumber); exists {
			fmt.Fprintf(w, `<Response>
					  <Dequeue phoneNumber="%s"/>
					</Response>`, virtualNumber)
			return
		}

		if active == "0" {
			// Toggle operator's availability status
			if sessionExists {
				writeOperators(op, false)
				writeSessions(sessionId, "")
			}
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
					fmt.Fprintf(w, `<Response><Say>Hello, please hold while we connect you to the next available operator</Say><Enqueue/></Response>`)
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

	go http.ListenAndServe(":"+port, nil)

	l, err := readline.New("\033[32m»\033[0m ")
	if err != nil {
		panic(err)
	}
	defer l.Close()

	for {
		line, err := l.Readline()
		if err != nil { // io.EOF
			break
		}

		tokens := strings.Split(line, " ")

		switch tokens[0] {
		case "operators":
			// Print active sessions
			operatorsMutex.Lock()
			fmt.Println(operators)
			operatorsMutex.Unlock()
		case "sessions":
			// Print active sessions
			sessionsMutex.Lock()
			fmt.Println(sessions)
			sessionsMutex.Unlock()
		case "transfer":
			// Transfer a session to a different operator
			if len(tokens) != 3 {
				fmt.Println("Usage: transfer <sessionId> <operatorNumber>")
			}else{
				transfer(tokens[1], tokens[2])
			}
		}		
		println(line)
	}
}
