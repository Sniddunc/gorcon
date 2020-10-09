# GoRCON

GoRCON is an RCON client implemented in Go. It implements Valve's Source RCON protocol documented [here](https://developer.valvesoftware.com/wiki/Source_RCON_Protocol).

# Usage

#### Go get it using:

`
go get github.com/sniddunc/gorcon
`

#### Initialize a client
```go
client, err := gorcon.NewClient(host, port, password)
if err != nil {
    // handle error
}
```

#### Attempt authentication
```go
err = client.Authenticate()
if err != nil {
    // handle error
}
```

#### Execute commands
```go
response, err := client.ExecCommand("help")
if err != nil {
    // handle error
}

fmt.Println(response)
```
#### Reconnect incase you get disconnected by the server
```go
err = client.Reconnect()
if err != nil {
    // handle error
}
```

