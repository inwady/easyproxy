package main

import (
    "flag"
    "net"
    "net/http"
    "fmt"
    "log"
    "crypto/tls"
    "strings"
)

func sendError(w http.ResponseWriter, message string) {
    fmt.Fprintf(w, "Error [%s]", message)
}

type proxy struct {
    baseAddr *net.TCPAddr
    host string

    ssl bool
    skipVerify bool
    redirectAccess bool
    v bool
}

func (p *proxy) normalizeServer(server string) string {
    if strings.Contains(server, ":") {
        return server
    }

    if !p.ssl {
        return server + ":80"
    } else {
        return server + ":443"
    }
}

func (p *proxy) connectWithServer(host string) (net.Conn, error) {
    var (
        connWithServer net.Conn
        err error
    )

    if !p.ssl {
        connWithServer, err = net.Dial("tcp", p.normalizeServer(host))
    } else {
        config := &tls.Config{
            InsecureSkipVerify: p.skipVerify,
        }

        connWithServer, err = tls.Dial("tcp", p.normalizeServer(host), config)
    }

    return connWithServer, err
}

func (p *proxy) connect() (net.Conn, error) {
    return p.connectWithServer(p.host)
}

func (p *proxy) verbose(logInfo interface{}) {
    if !p.v {
        return
    }

    log.Println(logInfo)
}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    bytesToServer, err := p.customRequest(r)
    if err != nil {
        sendError(w, "bad custom answer")
        return
    }

    p.verbose(string(bytesToServer))

    connWithServer, err := p.connect()
    if err != nil {
        sendError(w, "bad server addr")
        return
    }

    defer connWithServer.Close()

    _, err = connWithServer.Write(bytesToServer)
    if err != nil {
        sendError(w, err.Error())
        return
    }

    w.Header().Add("Proxy", "true")
    err = p.customResponse(r, w, connWithServer)
    log.Printf("success to %s", p.host)
}

func main() {
    server := flag.String("s", "", "server name")

    port := flag.String("p", "6556", "port number")
    ssl := flag.Bool("ssl", false, "ssl connect")
    skip := flag.Bool("k", false, "skip certificate")
    redirect := flag.Bool("r", false, "redirect access")
    verbose := flag.Bool("v", false, "verbose")
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
    proxyServer := &proxy{
        baseAddr: addr,
        host: *server,
        ssl: *ssl,
        skipVerify: *skip,
        redirectAccess: *redirect,
        v: *verbose,
    }

    proxyServer.initProxy(proxyServer)

    mux.Handle("/", proxyServer)

    baseServer := &http.Server{Addr: fmt.Sprintf(":%s", *port), Handler: mux}
    baseServer.SetKeepAlivesEnabled(false)
    err = baseServer.ListenAndServe()
    if err != nil {
        log.Panic(err)
    }
}