# Redis from scratch
A stripped down redis-server clone written as a hobby project. It is not intended to be used in production. No external dependencies were used for development.

## Prequisites
Have Go installed.

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
- Connect to the server with any redis client. For example, `redis-cli`. The server listens on port `6379` by default.
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
The `encoding/gob` package is used to asynchronously create a dump of the store in a binary file called dump.rdb every one second. This dump.rdb file is used to restore the data when then server is restarted in the same directory.
