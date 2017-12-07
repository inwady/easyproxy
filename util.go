package main

import (
    "net/http"
    "net"
    "bufio"
)

func readEasyHTTP(conn net.Conn) (*http.Response, error) {
    reader := bufio.NewReader(conn)
    resp, err := http.ReadResponse(reader, nil)
    if err != nil {
        return nil, err
    }

    return resp, nil
}
