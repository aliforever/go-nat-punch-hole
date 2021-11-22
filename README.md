# go-nat-punch-hole

This repository is to demonstrate how udp & tcp nat punchhole work.

## Workflow:
- There's a server with a map containing peers by their chosen name.
- Peer A & B register their ip addresses on the server with a chosen name.
- Both peers ask the server to see if there is another peer who has chosen that name.
- After getting each other's ips, they both start sending `Hello` messages to each other.

## How to run

### Server:
`go run main.go --type=server --local=0.0.0.0:8181`

### Peers:
`go run main.go --type=client --local=0.0.0.0:8182 --server=0.0.0.0:8181`

## This was tested with peers from Iran & Vancouver and was successful:
### Peer A:
```
peer_not_found
peer address is: X.X.X.X:8183
received MESSAGE:Hello from X.X.X.X:8183
received MESSAGE:Hello from X.X.X.X:8183
received MESSAGE:Hello from X.X.X.X:8183
received MESSAGE:Hello from X.X.X.X:8183
received MESSAGE:Hello from X.X.X.X:8183
received MESSAGE:Hello from X.X.X.X:8183
```

#### The `peer_not_found` errors is because peer A is the first peer alone on the server map.

### Peer B:
```
peer address is: Y.Y.Y.Y:8182
received MESSAGE:Hello from Y.Y.Y.Y:8182
received MESSAGE:Hello from Y.Y.Y.Y:8182
received MESSAGE:Hello from Y.Y.Y.Y:8182
received MESSAGE:Hello from Y.Y.Y.Y:8182
received MESSAGE:Hello from Y.Y.Y.Y:8182
received MESSAGE:Hello from Y.Y.Y.Y:8182
```
