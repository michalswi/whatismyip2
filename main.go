package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/gorilla/mux"
	"github.com/michalswi/whatismyip2/server"
)

var (
	snapshotLen int32 = 1024
	promiscuous bool  = false
	err         error
	timeout     time.Duration = 1 * time.Second
	handle      *pcap.Handle
)

var remoteIP string
var logger = log.New(os.Stdout, "IPer ", log.LstdFlags|log.Lshortfile|log.Ltime|log.LUTC)
var netInterface = flag.String("inter", "eth", "network interface listen to")

func main() {

	flag.Parse()
	serverPort := os.Getenv("SERVER_PORT")

	r := mux.NewRouter()
	srv := server.NewServer(r, serverPort)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("HTML Remote IP is: %s\n", remoteIP)
		var html = `
		<html>
		<p><b>IP</b>: %s</p>
		</html>
		`
		fmt.Fprintf(w, html, remoteIP)
	})

	r.HandleFunc("/ip", func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Raw Remote IP is: %s\n", remoteIP)
		fmt.Fprintf(w, "%s", remoteIP)
	})

	r.HandleFunc("/hz", func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Request from: %v\n", r.RemoteAddr)
		fmt.Fprintln(w, "ok")
	})

	// start server
	go func() {
		logger.Println("Server is ready to handle requests at port", serverPort)
		err := srv.ListenAndServe()
		if err != nil {
			logger.Fatalf("server failed to start: %v", err)
		}
	}()

	// start pcap
	go func() {
		logger.Println("Run packet capture..")
		var device = getInterfaceName(*netInterface)
		handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
		if err != nil {
			logger.Fatalf("pcap failed to start: %v", err)
		}
		defer handle.Close()
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			logger.Println("Inspecting packet...")
			remoteIP = getPacketInfo(packet)
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
		logger.Fatalf("could not gracefully shutdown the server: %v\n", err)
	}
	logger.Printf("Shutting down the server...\n")
	os.Exit(0)
}

func getInterfaceName(netInterface string) string {
	var device string
	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		if inter.Name != "lo" && strings.Contains(inter.Name, netInterface) {
			logger.Printf("Interface name: %s", inter.Name)
			device = inter.Name
		}
	}
	return device
}

func getPacketInfo(packet gopacket.Packet) string {
	var ipp string

	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
		logger.Printf("Source MAC: %s, Ethernet type: %s\n", ethernetPacket.SrcMAC, ethernetPacket.EthernetType)
	}

	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		if ip.TTL != 128 && ip.TTL != 64 {
			logger.Printf("SourceIP: %s, Protocol: %s, TTL: %d\n", ip.SrcIP, ip.Protocol, ip.TTL)
			ipp = fmt.Sprintf("%v", ip.SrcIP)
		}
	}

	if err := packet.ErrorLayer(); err != nil {
		logger.Println("error decoding some part of the packet:", err)
	}

	return ipp
}
