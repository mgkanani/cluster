#cluster
`cluster` is template for creating serverfarm. This is allowing to communicate in between them. One can either broadcast the message or send to individual server lying within `cluster`.


#Usage
## To install
```
go get github.com/mgkanani/cluster
go install github.com/mgkanani/cluster
```
## To run with default configuration:-
(go to cluster directory i.e. github.com/mgkanani/cluster and run below to see msgs' summary).
```
go test
```

#To customise parameters
###Modify config.json

To add more servers with different port and pids.
Pids must be in strictly order means 1,2,3,4 but not 1,2,4,7.
The order 1,4,3,2 will work perfectly.

###Modifying cluster_test.go

for Testing update values of variables total_servers count , and total_msgs.

for testing the msgs tranfer of size > 60000 bytes uncomment lines 48,57 and comment line 49.

#Default configurations:-
```
Total Servers :- 4
ipaddr:127.0.0.1 
ports :-12345,12346,12347,12348
Pids:-1,2,3,4
total_servers=4(linenum 18 in file cluster_test.go),
total_msgs=5001(linenum 19 in file cluster_test.go)
```


#References for basics of ZeroMQ and implemention
[http://nichol.as/zeromq-an-introduction](http://nichol.as/zeromq-an-introduction)

[http://zeromq.org/whitepapers:brokerless](http://zeromq.org/whitepapers:brokerless)

[https://github.com/alecthomas/gozmq](https://github.com/alecthomas/gozmq)

[http://stackoverflow.com/questions/14289256/cannot-convert-data-type-interface-to-type-string-need-type-assertion](http://stackoverflow.com/questions/14289256/cannot-convert-data-type-interface-to-type-string-need-type-assertion)
