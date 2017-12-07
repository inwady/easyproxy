package main

import (
    "net"
    "net/http"
    "fmt"
    "log"
    "io/ioutil"
)

func initProxy(proxy *proxy) {
    log.Println("init proxy")
}

func customRequest(r *http.Request) ([]byte, error) {
    b := []byte(fmt.Sprintf("GET / HTTP/1.1\r\nHost: mail.ru\r\n\r\n"))
    return b, nil
}

func customResponse(w http.ResponseWriter, connServer net.Conn) (error) {
    response, err := readEasyHTTP(connServer)
    if err != nil {
        return err
    }

    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return err
    }

    w.Write(body)
    return nil
}