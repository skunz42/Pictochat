package main;

import (
    "fmt"
    "os"
    "net"
    "bufio"
    "sync"
)

type connection struct {
    address string
    port string
    udp_or_tcp string
}

type client_connections struct {
    inputs []net.Conn
    receivers []net.Conn
}

var (
    mu sync.Mutex
    cli_conns client_connections
)

func make_connection(port string) *connection {
    c := connection{}
    hostname, err := os.Hostname()
    if err != nil {
        fmt.Println("Error fetching hostname")
        os.Exit(1)
    }
    ipaddr, err := net.LookupHost(hostname)
    if err != nil {
        fmt.Println("Error fetching IP address")
        os.Exit(1)
    }
    fmt.Println("Server address:", ipaddr[len(ipaddr)-1])
    c.address = ipaddr[len(ipaddr)-1]
    c.port = port
    c.udp_or_tcp = "tcp"
    return &c
}

func handle_message(conn net.Conn, user_connection *connection) {
    buffer, _ := bufio.NewReader(conn).ReadBytes('\n')

    if string(buffer) == "INPUT\n" {
        cli_conns.inputs = append(cli_conns.inputs, conn)
    } else if string(buffer) == "RECV\n" {
        cli_conns.receivers = append(cli_conns.receivers, conn)
    }

    fmt.Println("INPUTS:", len(cli_conns.inputs))
    fmt.Println("RECVS:", len(cli_conns.receivers))

    for {
        buffer, err := bufio.NewReader(conn).ReadBytes('\n')

        if err != nil {
            fmt.Println("Client left")
            conn.Close()
            return
        }

        fmt.Println("Message Received:", string(buffer))
    }
}

func start_server(my_conn *connection) {
    l, err := net.Listen(my_conn.udp_or_tcp, my_conn.address+":"+my_conn.port)
    if err != nil {
        fmt.Println("Error when listening:", err.Error())
        os.Exit(1)
    }

    defer l.Close()

    for {
        conn_obj, err := l.Accept()
        if err != nil {
            fmt.Println("Error on connection:", err.Error())
            return
        }

        fmt.Println("Now Entering:", conn_obj.RemoteAddr().String())
        go handle_message(conn_obj, my_conn)
    }
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Enter in the form ./server <port>")
        os.Exit(1)
    }

    port := os.Args[1]

    conn := make_connection(port)
    start_server(conn)
}
