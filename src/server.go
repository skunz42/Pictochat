package main;

import (
    "fmt"
    "os"
    "net"
    "bufio"
    "sync"
    "strings"
    "strconv"
)

type connection struct {
    address string
    port string
    udp_or_tcp string
    username string
}

type client_connections struct {
    inputs map[string]net.Conn
    receivers map[string]string
    num_messages uint64
    active_message string
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
        mu.Lock()
        cli_conns.num_messages += 1
        cli_conns.active_message = "Now entering: " + user_connection.username + "\n"
        mu.Unlock()
        for {
            buffer, err := bufio.NewReader(conn).ReadBytes('\n')

            if err != nil {
                fmt.Println("Goodbye,", user_connection.username)
                delete(cli_conns.inputs, user_connection.username)
                delete(cli_conns.receivers, user_connection.username)
                conn.Close()
                return
            } else if string(buffer) == ":list\n" {
                mu.Lock()
                fmt.Println("Users online:")
                cli_conns.num_messages += 1
                cli_conns.active_message = "Users online: "
                for k := range cli_conns.inputs {
                    fmt.Println(k)
                    cli_conns.active_message += (k[:len(k)-1] + ", ")
                }
                cli_conns.active_message += "("
                cli_conns.active_message += strconv.Itoa(len(cli_conns.inputs))
                cli_conns.active_message += "/8)\n"
                mu.Unlock()
            } else if string(buffer) == ":quit\n" {
                mu.Lock()
                cli_conns.num_messages += 1
                cli_conns.active_message = "Has exited the chat\n"
                mu.Unlock()
            } else {
                fmt.Println("Message Received:", string(buffer))
                /*for k,v := range cli_conns.receivers {
                    fmt.Fprintf(v, k+" TEST\n")
                }*/
                mu.Lock()
                cli_conns.num_messages += 1
                cli_conns.active_message = string(buffer)
                mu.Unlock()
            }
        }
    } else if string(string(buffer)[0]) == "R" {
        user_connection.username = parse_username(string(buffer))
        temp_split := strings.Split(string(buffer), ":")
        user_ip := temp_split[2]
        user_port := temp_split[3]
        user_both := user_ip+":"+user_port
        cli_conns.receivers[user_connection.username] = user_both[:len(user_both)-1]
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

func start_client() {
    size := cli_conns.num_messages
    for {
        mu.Lock()
        if size != cli_conns.num_messages {
            //fmt.Println("ACTIVE MSG:", cli_conns.active_message)
            for k,v := range cli_conns.receivers {
                temp_conn, err := net.Dial("tcp", v)
                if err != nil {
                    fmt.Println("Error connecting to recv:", err.Error())
                    os.Exit(1)
                }
                fmt.Fprintf(temp_conn, "<" + k + ">: " + cli_conns.active_message + "\n")
                temp_conn.Close()
            }
        }
        size = cli_conns.num_messages
        mu.Unlock()
    }
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Enter in the form ./server <port>")
        os.Exit(1)
    }

    port := os.Args[1]
    cli_conns.inputs = make(map[string]net.Conn)
    cli_conns.receivers = make(map[string]string)
    cli_conns.num_messages = 0

    conn := make_connection(port)
    go start_server(conn)
    go start_client()
    select {}
}
