package main;

import (
    "fmt"
    "os"
    "net"
    "bufio"
    "sync"
    "strings"
)

type connection struct {
    address string
    port string
    udp_or_tcp string
    username string
}

type client_connections struct {
    inputs map[string]net.Conn
    receivers map[string]net.Conn
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

func make_user_connection(conn net.Conn) *connection {
    c := connection{}
    ipport := conn.RemoteAddr().String()
    ipportarr := strings.Split(ipport, ":")
    c.address = ipportarr[0]
    c.port = ipportarr[1]
    c.udp_or_tcp = "tcp"
    return &c
}

func parse_username(message string) string {
    arr := strings.Split(message, ":")
    return arr[1]
}

func parse_user_connection(message string) *connection {
    arr := strings.Split(message, ":")
    c := connection{}
    c.address = arr[2]
    c.port = arr[3]
    c.udp_or_tcp = "tcp"
    c.username = arr[1]
    return &c
}

func handle_message(conn net.Conn) {

    user_connection := make_user_connection(conn)

    buffer, _ := bufio.NewReader(conn).ReadBytes('\n')

    if string(string(buffer)[0]) == "I" {
        user_connection.username = parse_username(string(buffer))
        cli_conns.inputs[user_connection.username] = conn
        fmt.Println("Now Entering:", user_connection.username)
        for {
            buffer, err := bufio.NewReader(conn).ReadBytes('\n')

            if err != nil {
                fmt.Println("Goodbye,", user_connection.username)
                delete(cli_conns.inputs, user_connection.username)
                //delete(cli_conns.receivers, user_connection.username)
                conn.Close()
                return
            } else if string(buffer) == ":list\n" {
                fmt.Println("Users online:")
                for k := range cli_conns.inputs {
                    fmt.Println(k)
                }
            } else {
                fmt.Println("Message Received:", string(buffer))
            }
        }
    } else if string(string(buffer)[0]) == "R" {
        user_connection.username = parse_username(string(buffer))
        cli_conns.receivers[user_connection.username] = conn
    }

    /*fmt.Println("INPUTS:", len(cli_conns.inputs))
    fmt.Println("RECVS:", len(cli_conns.receivers))*/

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

        go handle_message(conn_obj)
    }
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Enter in the form ./server <port>")
        os.Exit(1)
    }

    port := os.Args[1]
    cli_conns.inputs = make(map[string]net.Conn)
    cli_conns.receivers = make(map[string]net.Conn)

    conn := make_connection(port)
    start_server(conn)
}
