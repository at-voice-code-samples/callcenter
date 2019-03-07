##### This is a simple but functionally complete call center app
This version asks the caller to call back when all operators are on call  
See the "queue" branch for a version that places the caller on hold  
until an operator is free.

###### Build:
```
go build app.go
```

###### Run:
```
./app -port <port> -operators <comma-separated list of operators>
```

###### Example:
A callcenter with four operators, two with PSTN numbers and two with SIP addresses
```
./app -port 8080 -operators +2547xxxxxx,+2547xxxxxx,agent1@example.com,agent2@example.com
```
