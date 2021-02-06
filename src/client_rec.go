package main

import (
    "fmt"
    "os"
    "net"
)

type connection struct {
    address string
    port string
    udp_or_tcp string
    username string
}

func make_connection(ipaddr string, port string, username string) *connection {
    c := connection{}
    c.address = ipaddr
    c.port = port
    c.udp_or_tcp = "tcp"
    c.username = username
    return &c
}

func connect_to_server(my_conn * connection) {
    conn_obj, err := net.Dial(my_conn.udp_or_tcp, my_conn.address+":"+my_conn.port)
    if err != nil {
        fmt.Println("Error connecting to server:", err.Error())
        os.Exit(1)
    }

    fmt.Println("Sending initial message")
    fmt.Fprintf(conn_obj, "RECEIVER" + ":" + my_conn.username + "\n")
}

func main() {
    if len(os.Args) != 4 {
        fmt.Println("Enter in the form ./client <server port> <server ip address> <username>")
        os.Exit(1)
    }
    port := os.Args[1]
    ipaddr := os.Args[2]
    username := os.Args[3]
    conn := make_connection(ipaddr, port, username)
    connect_to_server(conn)
}
