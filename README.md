##### This is a simple but functionally complete call center app
This version is like the enqueue version, but with call forwarding  
capabilities. One can transfer a call in the CLI with the `transfer`  
command. Use the `operators` command to view operator statuses and
the `sessions` command to view sessions.

For this version, you need to supply your username and API key,
which are used to make the call transfer POST request

###### Build:
```
go build app.go
```

###### Run:
```
./app -port <port> -virtualNumber virtualNumber -username <name> -apikey <apikey> -operators <comma-separated list of operators>
```

###### Example:
A callcenter with four operators, two with PSTN numbers and two with SIP addresses
```
./app -port 8080 -virtualNumber +2547xxxxxxxx -username foo -apikey bar -operators +2547xxxxxx,+2547xxxxxx,agent1@example.com,agent2@example.com
```
