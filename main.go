package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	programName = "goserver"
	version     = "0.0.1"
)

var (
	network *string = flag.String("Network", "tcp4", "The network type to bind[tcp4|unix].")
	address *string = flag.String("Address", ":9000", "The bind address.")
)

var (
	g_listener net.Listener
	g_exit     bool
)

func process(request string) (response string) {
	// Todo: some logical codes.
	log.Println("request:", request)
	return `{"status":200}`
}

func read(c net.Conn) {
	defer c.Close()
	bw := bufio.NewWriter(c)
	br := bufio.NewReader(c)
	// Get request
	request, err := br.ReadString('\n')
	for err == nil {
		// Process
		response := process(request)
		if len(response) > 0 {
			// Response
			_, err = bw.WriteString(response)
			if err != nil {
				break
			}
			err = bw.Flush()
			if err != nil {
				break
			}
		}
		request, err = br.ReadString('\n')
	}
}

func listen() {
	for !g_exit {
		c, err := g_listener.Accept()
		if err != nil {
			if g_exit {
				break
			}
			// Rebind
			log.Println("Accept error:", err)
			time.Sleep(2 * time.Second)
			g_listener.Close()
			listener, _ := net.Listen(*network, *address)
			if listener != nil {
				g_listener = listener
			}
			continue
		}
		go read(c)
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s version[%s]\r\nUsage: %s [OPTIONS]\r\n", programName, version, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	var err error
	g_listener, err = net.Listen(*network, *address)
	if err != nil {
		log.Fatal("listen err:", err)
		os.Exit(-1)
	}

	go listen()
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	_ = <-ch
	g_exit = true
	g_listener.Close()
	fmt.Printf("Exit normally\n")
}
