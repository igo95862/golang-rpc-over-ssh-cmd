# Golang `net/rpc` over SSH using installed SSH program

This package implements a helper functions to launch an RPC client and server.
It uses the installed SSH program instead of native Go implementation from
`golang.org/x/crypto/ssh`. This has an advantage of being able to use SSH
configuration. (for example `${HOME}/.ssh/config`)

## How it works

Initiator launches the SSH command to connect to a remote server and launch a matching
binary. Matching binary reads from STDIN and writes to STDOUT. The SSH program will transparently
pipe STDOUT and STDIN between hosts.

STDERR from receiver will be shared with STDERR from client.

### Initiator

Initiator (also known as client) should call `sshrpc.NewSshRpcClient` to open a new RPC client.

`sshrpc.NewSshRpcClient` takes 2 arguments:

* `destination string` this is a host name that SSH will connect to. Example: `"myuser@somehost"`
* `arguments ...string` variable number of arguments that are required to launch the matching program.
  These arguments will be interpreted by the login shell on the remote host.

Returns has 2 returns:

* `*rpc.Client` see [net/rpc documentation](https://pkg.go.dev/net/rpc#Client).
  The client is already initialized and you can use `Call` and `Go` methods.
* `error` if some error happened.

### Receiver

Receiver (also known as server) should setup a `*rpc.Server` from [net/rpc](https://pkg.go.dev/net/rpc#Server)
and call `sshrpc.StartReceiving` method.

`sshrpc.StartReceiving` takes one argument which is the `*rpc.Server`. It will block indefinitely.

Receiver must NOT use stdin/stdout or the data stream will corrupt. Use STDERR if you want logging.

## Examples

See [examples folder](example/) for simple initiator and receiver implementation.

## TODO

* Better documentation
* Better README
* Add environment variable that controls SSH command line (like `GIT_SSH_COMMAND`)

