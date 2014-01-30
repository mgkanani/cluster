References for basics of ZeroMQ and implemention:-
http://zeromq.org/whitepapers:brokerless
http://nichol.as/zeromq-an-introduction -->this is best-ever among i have seen.
https://github.com/alecthomas/gozmq
http://stackoverflow.com/questions/14289256/cannot-convert-data-type-interface-to-type-string-need-type-assertion

Update the config.json file to add more servers(with different port and pids. Pids must be in strictly order means 1,2,3,4 but not 1,2,4,7. The order 1,4,3,2 will work perfectly.).

for Testing update total_servers count , and total_msgs values by modifying cluster_test.go file.

for testing the msgs tranfer of size > 60000 bytes uncomment lines 48,57 and comment line 49 of cluster_test.go

Default configurations:-
Total Servers :- 4 (ipadd:127.0.0.1 ;ports :-12345,12346,12347,12348;Pids:-1,2,3,4)
total_servers=4,total_msgs=10001

