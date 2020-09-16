package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/gorilla/mux"
	"github.com/michalswi/whatismyip/server"
)

// sudo SERVER_ADDR=8080 ./whatismyip

var (
	// device      string = "lo"
	snapshotLen int32 = 1024
	promiscuous bool  = false
	err         error
	timeout     time.Duration = 1 * time.Second
	handle      *pcap.Handle
)

var n = "none"

func main() {

	serverAddress := os.Getenv("SERVER_ADDR")

	r := mux.NewRouter()
	srv := server.NewServer(r, serverAddress)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Remote IP is: %s\n", n)
		fmt.Fprintf(w, "%s", n)
	})

	r.HandleFunc("/hz", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request from: %v\n", r.RemoteAddr)
		fmt.Fprintln(w, "ok")
	})

	// in docker initial request to localhost/lo is needed
	r.HandleFunc("/in", func(w http.ResponseWriter, r *http.Request) {
		_, err := http.Get(fmt.Sprintf("http://localhost:%s", serverAddress))
		if err != nil {
			log.Printf("Initial request: %v\n", err)
		}
	})

	// start server
	go func() {
		log.Printf("Starting server on port %s..", serverAddress)
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// start pcap
	go func() {
		log.Println("Run pcap..")
		var device = getInterfaceName()
		handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
		if err != nil {
			log.Fatal(err)
		}
		defer handle.Close()
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			log.Println("Inspecting packet...")
			n = getPacketInfo(packet)
		}
	}()

	// shutdown server
	gracefulShutdown(srv)
}

func gracefulShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interruptChan
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	log.Printf("Shutting down the server...\n")
	os.Exit(0)
}

func getInterfaceName() string {
	var device string
	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		if inter.Name != "lo" {
			// # todo
			fmt.Println("Name  :", inter.Name)
			device = inter.Name
		}
	}
	fmt.Println(device)
	return device
}

func getPacketInfo(packet gopacket.Packet) string {
	var ipp string

	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
		log.Printf("Source MAC: %s, Ethernet type: %s\n", ethernetPacket.SrcMAC, ethernetPacket.EthernetType)
	}

	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		log.Printf("SourceIP: %s, Protocol: %s\n", ip.SrcIP, ip.Protocol)
		ipp = fmt.Sprintf("%v", ip.SrcIP)
	}
	if err := packet.ErrorLayer(); err != nil {
		log.Println("Error decoding some part of the packet:", err)
	}
	return ipp
}
