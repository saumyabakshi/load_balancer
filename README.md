# load_balancer
Building a load balancer in Golang

For the coding challenge - <https://codingchallenges.fyi/challenges/challenge-load-balancer> , I decided to build a simple load balancer in Golang. The load balancer supports round-robin and least connections strategy to distribute incoming requests to a list of servers. 

## How to run
```bash
go run main.go
```
with flags:
```bash
go run main.go -port=8080 -strategy=least-connections
```
port: The port number to run the load balancer on. Default is 8000.
strategy: The strategy to use for load balancing. Default is round-robin.

It implements health check for the servers using a goroutine and gets the active connections for each server.


#References:
- https://betterprogramming.pub/building-a-load-balancer-in-go-3da3c7c46f30
- https://github.com/swayne275/load-balancer-proxy/tree/main
- https://codingchallenges.fyi/challenges/challenge-load-balancer

