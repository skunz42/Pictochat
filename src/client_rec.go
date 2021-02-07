package main

import (
    "fmt"
    "os"
    "os/exec"
    "net"
    "bufio"
    "strings"
)

type connection struct {
    address string
    port string
    udp_or_tcp string
    username string
}

var(
    local_ipp string
    local_ipr []string
    local_ip string
)

func make_connection(ipaddr string, port string, username string) *connection {
    c := connection{}
    c.address = ipaddr
    c.port = port
    c.udp_or_tcp = "tcp"
    c.username = username
    return &c
}

func connect_to_server(my_conn *connection, user_port string) {
    l, err := net.Listen("tcp", local_ip+":"+user_port)
    if err != nil {
        fmt.Println("Error when listening:", err.Error())
        os.Exit(1)
    }

    defer l.Close()

    fmt.Println("Welcome to Pictochat!")

    for {
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error on connection:", err.Error())
            return
        }
        buffer, _ := bufio.NewReader(conn).ReadBytes('\n')

        if (string(buffer) == ":quit\n") {
            break
        } else if (len(string(buffer)) > 5 && string(buffer)[:5] == ":draw") {
            temp := string(buffer)[5:]
            //fmt.Print(temp)
            var k int = len(temp)/120
            for i := 0; i <= k; i++ {
                if (i+1)*120 < len(temp) {
                    fmt.Println(temp[i*120:(i+1)*120])
                } else {
                    fmt.Println(temp[i*120:len(temp)])
                }
            }
        } else {
            fmt.Print(string(buffer))
        }
    }
}

func send_msg(my_conn * connection, user_port string) {
    conn_obj, err := net.Dial(my_conn.udp_or_tcp, my_conn.address+":"+my_conn.port)
    local_ipp = conn_obj.LocalAddr().String()
    local_ipr = strings.Split(local_ipp, ":")
    local_ip = local_ipr[0]
    //local_port := local_ipr[1]
    if err != nil {
        fmt.Println("Error connecting to server:", err.Error())
        os.Exit(1)
    }

    //fmt.Println("Sending initial message")
    fmt.Fprintf(conn_obj, "RECEIVER" + ":" + my_conn.username + ":" + local_ip + ":" + user_port + "\n")
}

func main() {
    c := exec.Command("clear")
    c.Stdout = os.Stdout
    c.Run()
    if len(os.Args) != 5 {
        fmt.Println("Enter in the form ./client <server port> <server ip address> <username> <user_port>")
        os.Exit(1)
    }
    port := os.Args[1]
    ipaddr := os.Args[2]
    username := os.Args[3]
    user_port := os.Args[4]
    conn := make_connection(ipaddr, port, username)
    send_msg(conn, user_port)
    connect_to_server(conn, user_port)
}
