// Implementation of a MultiEchoServer. Students should write their code in this file.

package project1

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
)

var (
	msgBufferLimit      = 100
	debug_num_msgs_read = 0
	increment_ID        = 0
)

type client struct {
	id           int
	conn         net.Conn
	clientBuffer (chan string)
	closed       (chan int)
}

type multiEchoServer struct {
	listener     net.Listener
	serverBuffer (chan string)
	clientsMap   map[int]client
}

func NewClient(conn net.Conn, clientId int) client {
	current_client := client{
		conn:         conn,
		closed:       make(chan int, 1), // 1 in the channel means closed
		clientBuffer: make(chan string, msgBufferLimit),
		id:           clientId,
	}
	return current_client
}

// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
	current := multiEchoServer{
		serverBuffer: make(chan string),
		clientsMap:   make(map[int]client),
	}

	// current := new(multiEchoServer)
	// current.listener = nil
	// current.messagesCh = make(chan string)
	// current.connMsgMap = make(map[net.Conn](chan string))

	// mu := &Mutex{make(chan int, 1)}
	// current := &multiEchoServer{messagesCh: make(chan string),
	// 			connMsgMap: make(map[net.Conn](chan string, msgBufferLimit))}
	return &current
}

func (mes *multiEchoServer) Start(port int) error {

	// Start the server on listening to given port
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	Error("Error getting server up: %v\n", err)
	// fmt.Printf("server is up  \n")
	mes.listener = listener

	go acceptConnections(mes)
	return nil
}

func (mes *multiEchoServer) Close() {
	for clientID := range mes.clientsMap {
		mes.clientsMap[clientID].closed <- 1
		// currentLen := len(mes.clientsMap[clientID].closed)
		// fmt.Printf("attempting to close the client, client chan len is %v\n",
		// currentLen)
		connection := mes.clientsMap[clientID].conn
		connection.Close()
	}
	mes.listener.Close()
}

func (mes *multiEchoServer) Count() int {
	// fmt.Printf("count is %v\n", len(mes.clientsMap))
	return len(mes.clientsMap)
}

// TODO: add additional methods/functions below!
func Error(err_msg string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, err_msg, err)
		os.Exit(1)
	}
}

func acceptConnections(mes *multiEchoServer) {
	// Just keep running
	go broadcast(mes)

	for {

		// accept connections
		conn, err := mes.listener.Accept()
		// Error("Error accepting connection: %v\n", err)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accepting connection: %v\n", err)
			return
		}
		// fmt.Printf("server is accepting connections \n")

		// once a connection is received, we add it to our list
		clientID := increment_ID
		current_client := NewClient(conn, clientID)
		mes.clientsMap[increment_ID] = current_client
		increment_ID++

		// fmt.Printf("adding new client %v\n", clientID)
		go readMsg(conn, mes, clientID)
		// fmt.Printf("POOOOOPPPPPP")

		go writeMsg(clientID, mes)
	}
}

// read in the message from the current connection
func readMsg(conn net.Conn, mes *multiEchoServer, clientID int) {
	reader := bufio.NewReader(conn)

	for {
		select {
		case <-mes.clientsMap[clientID].closed:
			// fmt.Printf("server is closed \n")
			// fmt.Printf("deleting client in readMsg \n")
			delete(mes.clientsMap, clientID)
		default:
			message, err := reader.ReadString('\n')
			// If read isn't possible, conn is closed
			if err != nil {
				// fmt.Printf("deleting client in readMsg error \n")
				delete(mes.clientsMap, clientID)
				return
			}
			// debug_num_msgs_read++
			// fmt.Printf("can read from connection: %v\n", clientID)

			mes.serverBuffer <- message
		}
	}
}

func broadcast(mes *multiEchoServer) {
	for {
		currentMsg := <-mes.serverBuffer
		for clientID := range mes.clientsMap {
			// add message to each client's channel
			current_client_buffer := mes.clientsMap[clientID].clientBuffer
			if len(current_client_buffer) < msgBufferLimit {
				// fmt.Printf("adding to the client's buffer ")
				mes.clientsMap[clientID].clientBuffer <- currentMsg
			}
		}
	}
}

// loop through all the clients
// and write out a message for each client
func writeMsg(clientID int, mes *multiEchoServer) {
	for {
		select {
		case <-mes.clientsMap[clientID].closed:
			// fmt.Printf("server is closed \n")
			delete(mes.clientsMap, clientID)
			// fmt.Printf("deleting client in writeMsg \n")
		default:
			// for clientID := range mes.clientsMap {
			currentMsg := <-mes.clientsMap[clientID].clientBuffer
			// fmt.Printf("i am writing ")
			fmt.Fprintf(mes.clientsMap[clientID].conn, "%v", currentMsg)
			// }
		}
	}
}

// func readMsg(conn net.Conn) {
// 	// create reader from connection
// 	reader := bufio.NewReader(conn)
//
// 	for {
// 		// read message
// 		message, err := reader.ReadString('\n')
// 		Error(err, "Error reading from connection: %v\n")
//
//
// 	}
//
// }
