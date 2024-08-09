# cjdns-peers-finder

This is a server for retrieving "good" nodes for peering. 
On each request the information about the cjdns node is stored and combining peering information from cjdns route-server it returns up to 3 "good" nodes for peering in random, so that each node gets different peers.


## Features
* API Endpoint: Provides an endpoint to retrieve peer information in JSON format.
* Periodic Peer Testing: Runs peer tests every hour using a goroutine.
* Logging: Logs important events and errors for monitoring and debugging.
## Prerequisites
* Go 1.16 or higher
* CJDNS installed and configured
* Environment variable CJDNS_PATH set to the path of CJDNS tools

## Installation
Clone the repository:

```git clone <repository-url>cd <repository-directory>```

Build the project:

```go build -o cjdns-peers-finder```

Usage
Run the server:

```CJDNS_PATH=/home/user/cjdns/tools ./cjdns-peers-finder```

Access the API endpoint:

```curl -X POST -H "Content-Type: application/json" -d '{"name":"Server name","login":"default","pasword":"pwd","ip":"192.168.1.1","port":9999,"publicKey":"026whhvh3j3bnmv8vxhwzjf6b31nt0tr4kmqv67bxqutnhufxz00.k"}' http://localhost:8090/api/peers```


## Logging
The server logs important events and errors. Ensure that the logger is properly configured to capture these logs.