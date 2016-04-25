package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
)

func main() {
	ip := getLocalIP("0.0.0.0")
	port := getPort("8000")

	http.Handle("/", http.FileServer(http.Dir("./")))

	fmt.Printf("Listening on %s:%s\n", ip, port)
	listenString := fmt.Sprintf(":%s", port)
	http.ListenAndServe(listenString, logHandler(http.DefaultServeMux))
}

func logHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lw := &loggingResponseWriter{w, 0}
		handler.ServeHTTP(lw, r)
		if lw.status == 0 {
			lw.status = 200
		}
		message := fmt.Sprintf("[%s] %d %s %s", stripPort(r.RemoteAddr), lw.status, r.Method, r.URL)

		if lw.status <= 300 {
			color.Green(message)
		} else if lw.status <= 400 {
			color.Cyan(message)
		} else {
			color.Red(message)
		}
	})
}

func stripPort(addr string) string {
	return strings.Split(addr, ":")[0]
}

type loggingResponseWriter struct {
	ResponseWriter http.ResponseWriter
	status         int
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

func (w loggingResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func getLocalIP(fallback string) string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return fallback
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return fallback
}

func getPort(fallback string) string {
	port := os.Getenv("port")
	if len(port) == 0 {
		port = os.Getenv("PORT")
	}
	if len(port) == 0 {
		port = fallback
	}
	return port
}
