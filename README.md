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
./app -port 8080 -virtualNumber +2547xxxxxxxx -operators +2547xxxxxx,+2547xxxxxx,agent1@example.com,agent2@example.com
```

###### Call queueing and dequeueing
For this version, when all available operators are busy, consequent callers are placed on hold.  
To dequeue callers on hold, a free operator simply calls the virtualnumber and gets connected  
to one of the callers on hold.  
Note that for operators on SIP phones, the callback for each SIP address has to be configured to
respond with the dequeue action for this to work.