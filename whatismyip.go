package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	device      string = "lo"
	snapshotLen int32  = 1024
	promiscuous bool   = false
	err         error
	timeout     time.Duration = 1 * time.Second
	handle      *pcap.Handle
)

func main() {

	serverAddress := os.Getenv("SERVER_ADDR")

	http.HandleFunc("/", getIP)

	http.HandleFunc("/hz", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request from: %v\n", r.RemoteAddr)
		fmt.Fprintln(w, "ok")
	})

	fmt.Printf("Starting server on port %s...\n", serverAddress)
	port := fmt.Sprintf(":%s", serverAddress)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("error, %s\n", err)
	}
}

func getIP(w http.ResponseWriter, r *http.Request) {
	var n string
	handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		n = getPacketInfo(packet)
		break
	}
	log.Printf("Remote IP is: %s\n", n)
	fmt.Fprintf(w, "%s", n)
}

func getPacketInfo(packet gopacket.Packet) string {
	var ipp string
	// ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	// if ethernetLayer != nil {
	// 	ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
	// 	log.Printf("Source MAC: %s, Ethernet type: %s\n", ethernetPacket.SrcMAC, ethernetPacket.EthernetType)
	// }
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
