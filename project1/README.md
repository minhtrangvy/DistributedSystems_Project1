# Project 1: Multi-client echo server

This subrepository contains the instructions and the starter code that you will use as the basis of your multi-client echo server implementation. It also contains the tests that we will use to test your implementation, and an example 'server runner' binary that you might find useful for your own testing purposes.

The goal of this assignment is to get up to speed on the Go programming language and to help remind you about the basics of socket programming. In this assignment you will implement a simple multi-client echo server in Go: every message sent by a client should be echoed to all connected clients.

## Server characteristics

Your multi-client echo server must have the following characteristics:

1. The server must manage and interact with its clients concurrently using goroutines and channels. Multiple clients should be able to connect/disconnect to the server simultaneously.

2. The server should assume that all messages are line-oriented, separated by newline (`\n`) characters. When the server reads a newline-terminated message from a client, it must respond by writing that exact message (up to and _including_ the newline character) to all connected clients, including the client that sent the message.

3. The server must be responsive to slow-reading clients. To better understand what this means, consider a scenario in which a client does not call Read for an extended period of time. If during this time the server continues to write messages to the client’s TCP connection, eventually the TCP connection’s output buffer will reach maximum capacity and subsequent calls to `Write` made by the server will block.
To handle these cases, your server should keep a queue of at most 100 outgoing messages to be written to the client at a later time. Messages sent to a slow-reading client whose outgoing message buffer has reached the maximum capacity of 100 should simply be dropped. If the slow-reading client starts reading again later on, the server should make sure to write any buffered messages in its queue back to the client. (Hint: use a buffered channel to implement this property).

## Requirements

This project is intentionally open-ended and there are many possible solutions. That said, your implementation must meet the following four requirements:

1. The project must be done individually. You are not allowed to use any code that you have not written yourself.

2. Your code may not use locks and mutexes. All synchronization must be done using goroutines, channels, and Go’s channel-based select statement (not to be confused with the low-level socket select that you might use in C, which is also not allowed).

3. You may only use the following packages: `bufio`, `fmt`, `net`, `os`, and `strconv`.

4. You must format your code using `go fmt` and must follow Go’s standard naming conventions. See the [Formatting](https://golang.org/doc/effective_go.html#formatting) and [Names](https://golang.org/doc/effective_go.html#names) sections of Effective Go for details.

We don’t expect your solutions to be overly-complicated. As a reference, our sample solution is a little over 100 lines including sparse comments and whitespace. We do, however, _highly recommend_ that you familiarize yourself with Go’s concurrency model before beginning the assignment.

## The starter code

To download the code for this class, follow [these instructions](https://github.com/jnylam/cs189a).

The starter code for this project, which is under `$GOPATH/src/github.com/jnylam/cs189a/project1`, consists of the following files:

1. `server_impl.go` is the only file you should modify, and is where you will add your code as you implement your multi-client echo server.

2. `server_api.go` contains the interface and documentation for the MultiEchoServer you will be implementing in this project. You should not modify this file.

3. `server_test.go` contains the tests that we will run to grade your submission.

If at any point you have any trouble with building, installing, or testing your code, the article
titled [How to Write Go Code](http://golang.org/doc/code.html) is a great resource for understanding
how Go workspaces are built and organized. You might also find the documentation for the
[`go` command](http://golang.org/cmd/go/) to be helpful. As always, feel free to post your questions
on Piazza as well.

## Running the official tests

To test your submission, we will execute the following command from inside the
`src/github.com/jnylam/cs189a/project1` directory:

```sh
$ go test
```

We will also check your code for race conditions using Go's race detector by executing
the following command:

```sh
$ go test -race
```

To execute a single unit test, you can use the `-test.run` flag and specify a regular expression
identifying the name of the test to run. For example,

```sh
$ go test -race -test.run TestBasic1
```

## Testing your implementation using `srunner`

To make testing your server a bit easier (especially during the early stages of your implementation
when your server is largely incomplete), we have given you a simple `srunner` (server runner)
program that you can use to create and start an instance of your `MultiEchoServer`. The program
simply creates an instance of your server, starts it on a default port, and blocks forever,
running your server in the background.

To compile and build the `srunner` program into a binary that you can run, execute the three
commands below:

```bash
$ go install github.com/jnylam/cs189a/project1/srunner
$ $GOPATH/bin/srunner
```

The `srunner` program won't be of much use to you without any clients. It might be a good exercise
to implement your own `crunner` (client runner) program that you can use to connect with and send
messages to your server. We have provided you with an unimplemented `crunner` program that you may
use for this purpose if you wish. Whether or not you decide to implement a `crunner` program will not
affect your grade for this project.

You could also test your server using Netcat (i.e. run the `srunner`
binary in the background, execute `nc localhost 9999`, type the message you wish to send, and then
click enter).

## Acknowledgements

This assignment was given as part of [Distributed Systems](http://www.cs.cmu.edu/~dga/15-440/S14/index.html) taught by David Andersen and Srini Seshan.
