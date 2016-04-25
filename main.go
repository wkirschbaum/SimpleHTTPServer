package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

func main() {
	ip := getLocalIP()
	port := os.Getenv("port")
	if len(port) == 0 {
		port = "8000"
	}

	http.Handle("/", http.FileServer(http.Dir("./")))

	fmt.Printf("Listening on %s:%s\n", ip, port)
	http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), nil)
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
