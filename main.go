package main

import (
    "flag"
    "net"
    "net/http"
    "fmt"
    "log"
)

func sendError(w http.ResponseWriter, message string) {
    fmt.Fprintf(w, "Error [%s]", message)
}

type proxy struct {
    baseAddr *net.TCPAddr
}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    bytesToServer, err := customRequest(r)
    if err != nil {
        sendError(w, "bad custom answer")
        return
    }

    tcpConn, err := net.DialTCP("tcp", nil, p.baseAddr)
    if err != nil {
        sendError(w, "bad server addr")
        return
    }

    _, err = tcpConn.Write(bytesToServer)
    if err != nil {
        sendError(w, err.Error())
    }

    err = customResponse(w, tcpConn)
    log.Println("success!")

    tcpConn.Close()
}

func main() {
    server := flag.String("s", "", "server name")
    port := flag.String("p", "6556", "port number")
    flag.Parse()

    if *server == "" {
        flag.Usage()
        return
    }

    addr, err := net.ResolveTCPAddr("tcp", *server)
    if err != nil {
        log.Panic(err)
    }

    mux := http.NewServeMux()
    proxyServer := &proxy{baseAddr: addr}
    initProxy(proxyServer)

    mux.Handle("/", proxyServer)

    baseServer := &http.Server{Addr: fmt.Sprintf(":%s", *port), Handler: mux}
    baseServer.SetKeepAlivesEnabled(false)
    err = baseServer.ListenAndServe()
    if err != nil {
        log.Panic(err)
    }
}