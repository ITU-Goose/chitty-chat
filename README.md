# Chitty chat
## Run the client directly
Run:
```
go run cmd/client/main.go
```
Arguments:
- `--host {ip}` specify host of server to connect to
- `--port {port}` specify port of server 
- `--random` constantly send a message with a random number to server  

## Run the serve
Run:
```
go run cmd/server/main.go
```

## Deploy
Run:
```
docker-compose up --scale client=3
```