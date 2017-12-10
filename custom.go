package main

import (
    "net/http"
    "log"
    "io"
)

func (p *proxy) initProxy(proxy *proxy) {
    log.Println("init proxy")
}

func (p *proxy) customRequest(r *http.Request) ([]byte, error) {
    return parseEasyHTTPWithHost(r, p.host)
}

func (p *proxy) customResponse(r *http.Request, w http.ResponseWriter, connServer io.Reader) (error) {
    response, err := readEasyHTTP(connServer)

    if p.redirectAccess &&
        (response.StatusCode == http.StatusMovedPermanently ||
         response.StatusCode == http.StatusFound) {
        /* try 1 time */
        location := response.Header.Get("Location")
        r.URL, err = r.URL.Parse(location)
        if err != nil {
            return err
        }

        r.Host = r.URL.Host

        bytes, err := parseEasyHTTP(r)
        if err != nil {
            return err
        }

        conn, err := p.connectWithServer(r.Host)
        if err != nil {
            return err
        }

        defer conn.Close()

        conn.Write(bytes)
        response, err = readEasyHTTP(conn)
    }

    if err != nil {
        return err
    }

    return copyEasyHTTP(w, response)
}