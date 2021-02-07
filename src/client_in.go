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
    username string
}

func make_connection(port string, ipaddr string, username string) *connection {
    c := connection{}
    c.address = ipaddr
    c.port = port
    c.udp_or_tcp = "tcp"
    c.username = username
    return &c
}

func connect_to_server(my_conn *connection) {
    conn_obj, err := net.Dial(my_conn.udp_or_tcp, my_conn.address+":"+my_conn.port)
    if err != nil {
        fmt.Println("Error connecting to server:", err.Error())
        os.Exit(1)
    }

    //fmt.Println("Sending initial message")
    fmt.Fprintf(conn_obj, "INPUT" + ":" + my_conn.username + "\n")
    //fmt.Println("Connected to the server successfully!")

    for {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("Text to send: ")
        text, _ := reader.ReadString('\n')
        fmt.Fprintf(conn_obj, text + "\n")
        c := exec.Command("clear")
        c.Stdout = os.Stdout
        c.Run()
        if string(text) == ":quit\n" {
            //fmt.Println("quitting lol")
            return
        }
    }
}

func main() {
    c := exec.Command("clear")
    c.Stdout = os.Stdout
    c.Run()
    if len(os.Args) != 4 {
        fmt.Println("Enter in the form ./client <server port> <server ip> <username>")
        os.Exit(1)
    }
    port := os.Args[1]
    ipaddr := os.Args[2]
    username := os.Args[3]
    conn := make_connection(port, ipaddr, username)
    connect_to_server(conn)
}
