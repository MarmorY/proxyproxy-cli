package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/Neothorn23/proxyproxy"
	"github.com/Neothorn23/proxyproxy-handler/sspi"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

var (
	version string
)

func main() {
	fmt.Printf("proxyproxy-cli %s\n\n", version)

	log.SetHandler(cli.New(os.Stdout))

	//Define CLI flags
	destinationProxy := flag.String("proxy", "", "destination proxy: <ip addr>:<port>")
	listenAddress := flag.String("listen", "127.0.0.1:3128", "adress to list on: [ip addr]:<port>")
	debug := flag.Bool("v", false, "Verbose output")

	flag.Parse()

	if *destinationProxy == "" {
		fmt.Println("Parameter \"proxy\" is not set.")
		flag.PrintDefaults()
		return
	}

	if *listenAddress == "" {
		fmt.Println("Parameter \"listen\" is not set.")
		flag.PrintDefaults()
		return
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("Verbos output is enabled.")
	}

	log.Infof("Listening on %v", *listenAddress)
	log.Infof("Connection to %v", *destinationProxy)

	ln, err := net.Listen("tcp", *listenAddress)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	authHandler, err := sspi.NewSSPIAuthHandler()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		go handleConnecion(conn, *destinationProxy, authHandler)
	}
}

func handleConnecion(clientConn net.Conn, proxyServer string, authHandler proxyproxy.NtlmAuhtHandler) {

	logger := log.Log.(*log.Logger)
	cliEventLogger := NewCliEventLogger(log.NewEntry(logger))

	proxyConn, err := net.DialTimeout("tcp", proxyServer, 10*time.Second)
	if err != nil {
		log.Fatalf("Error opening connection to proxy: %v", err)
	}

	communication, err := proxyproxy.NewProxyCommunication(clientConn, proxyConn, authHandler, cliEventLogger)
	if err != nil {
		log.Errorf("Error creating a new proxy communication: %v", err)
	} else {
		communication.HandleConnection()
	}

}
