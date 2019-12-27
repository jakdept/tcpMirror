package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/alecthomas/kingpin"
)

var (
	protocol = kingpin.Flag("protocol", "protocol type (tcp, tcp4, tcp6)").Default("tcp4").String()
	listen   = kingpin.Flag("listen", "listen address").Default(":7070").String()
	announce = kingpin.Flag("announce", "initial message").Default("will mirror everything back\n").String()
	timeout  = kingpin.Flag("timeout", "idle connection timeout").Default("10s").Duration()
)

func main() {
	kingpin.Parse()
	l, err := net.Listen(*protocol, *listen)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		go Handle(conn)
	}
}

func Handle(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("got a connection from %s\n", conn.RemoteAddr())
	defer func() {
		fmt.Printf("closed connection from %s\n", conn.RemoteAddr())
	}()

	_, err := conn.Write([]byte(*announce))
	if err != nil {
		log.Println(err)
	}

	_, err = io.Copy(conn, conn)
	if err != nil {
		log.Println(err)
		return
	}
}
