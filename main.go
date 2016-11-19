//In the name of God
// Author: Shayan Hosseini
// with major help of Kianoosh

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

// Server information ...
type Server struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
}

func getServerInfo(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Could not read hostname"
	}
	ip, err := externalIP()
	if err != nil {
		ip = "Could not find IP address"
	}

	resp := Server{
		Hostname: hostname,
		IP:       ip,
	}

	jsonResp, err := json.Marshal(&resp)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		fmt.Println(string(jsonResp))
		w.Write(jsonResp)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getServerInfo)

	log.Fatalln(http.ListenAndServe(":8080", mux))
}
