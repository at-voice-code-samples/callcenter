##### This is a simple but functionally complete call center app

###### Build:
```
go build app.go
```

###### Run:
```
./app -port <port> -virtualNumber virtualNumber -operators <comma-separated list of operators>
```

###### Example:
A callcenter with four operators, two with PSTN numbers and two with SIP addresses
```
./app -port 8080 -operators +2547xxxxxx,+2547xxxxxx,agent1@example.com,agent2@example.com
```

###### Call queueing and dequeueing
For this version, when all available operators are busy, consequent caller are placed on hold.
To dequeue callers on hold, a free operator simply calls the virtualnumber and gets connected
to one of the callers on hold