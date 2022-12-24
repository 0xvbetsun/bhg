package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	ports  *int
	target *string
)

func init() {
	ports = flag.Int("ports", 1024, "# of ports to scan")
	target = flag.String("target", "", "target ip")
}

func main() {
	flag.Parse()

	if *target == "" {
		log.Fatal("target is not specified")
	}
	var wg sync.WaitGroup

	log.Printf("started scanning of %q, from 1 to %d ports", *target, *ports)

	for port := 1; port <= *ports; port++ {
		wg.Add(1)
		go sanTCPConn(port, &wg)
	}
	wg.Wait()
}

func sanTCPConn(port int, wg *sync.WaitGroup) {
	defer wg.Done()

	address := fmt.Sprintf("%s:%d", *target, port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return
	}
	defer conn.Close()

	if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		log.Printf("conn set deadline: %v", err)
		return
	}

	var sb strings.Builder
	buf := make([]byte, 256)

	n, err := conn.Read(buf)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return
		}
		log.Printf("conn read: %v", err)
		return
	}
	sb.Write(buf[:n])
	log.Printf("addr: %s:%d - %s", *target, port, sb.String())
}
