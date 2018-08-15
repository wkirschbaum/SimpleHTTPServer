package main

import (
	"flag"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
  "log"

	"github.com/fatih/color"
	"github.com/wkirschbaum/localip"
)

func main() {
	var port string
	var directory string
	var help bool
	flag.BoolVar(&help, "h", false, "help")
	flag.StringVar(&port, "p", "8000", "port")
	flag.StringVar(&directory, "d", "./", "directory")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		return
	}

	ip := localip.GetLocalIP("0.0.0.0")

	http.Handle("/", http.FileServer(http.Dir(directory)))

	absPath, _ := filepath.Abs(directory)
	fmt.Printf("Serving %s\n", absPath)
	fmt.Printf("Listening on %s:%s\n", ip, port)
	listenString := fmt.Sprintf(":%s", port)
  log.Fatal( http.ListenAndServe(listenString, logHandler(headerHandler(http.DefaultServeMux))))
}

func headerHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
		handler.ServeHTTP(w, r)
  })
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
