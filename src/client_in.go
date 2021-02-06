package main

import (
    "fmt"
    "os"
    "os/exec"
    "net"
    "bufio"
)

type connection struct {
    address string
    port string
    udp_or_tcp string
}

func make_connection(port string, ipaddr string) *connection {
    c := connection{}
    c.address = ipaddr
    c.port = port
    c.udp_or_tcp = "tcp"
    return &c
}

func connect_to_server(my_conn *connection) {
    conn_obj, err := net.Dial(my_conn.udp_or_tcp, my_conn.address+":"+my_conn.port)
    if err != nil {
        fmt.Println("Error connecting to server:", err.Error())
        os.Exit(1)
    }

    fmt.Println("Sending initial message")
    fmt.Fprintf(conn_obj, "INPUT" + "\n")
    fmt.Println("Connected to the server successfully!")

    for {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("Text to send: ")
        text, _ := reader.ReadString('\n')
        fmt.Fprintf(conn_obj, text + "\n")
        if string(text) == ":quit\n" {
            fmt.Println("quitting lol")
            return
        }
    }
}

func main() {
    c := exec.Command("clear")
    c.Stdout = os.Stdout
    c.Run()
    if len(os.Args) != 3 {
        fmt.Println("Enter in the form ./client <server port> <server ip>")
        os.Exit(1)
    }
    port := os.Args[1]
    ipaddr := os.Args[2]
    conn := make_connection(port, ipaddr)
    connect_to_server(conn)
}
