# Redis from scratch
A stripped down redis-server clone written as a hobby project. It is not intended to be used in production. No external dependencies were used for development.

## Usage
- Clone the repo
- Run by executing 
```
go run .
```
- Or build and run using 
```
go build . && ./redis
```
- Connect to the server with any redis client. For example, `redis-cli`
```
$ redis-cli
```

## Supported commands
- PING
- SET
- GET
- DEL
- MSET
- MGET
- HSET
- HGET
- HDEL
- EXISTS
- HEXISTS

## Persistence
At the moment, persistence isn't implemented.
