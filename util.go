package main

import (
    "net/http"
    "bufio"
    "net/http/httputil"
    "io/ioutil"
    "io"
)

func readEasyHTTP(conn io.Reader) (*http.Response, error) {
    reader := bufio.NewReader(conn)
    resp, err := http.ReadResponse(reader, nil)
    if err != nil {
        return nil, err
    }

    return resp, nil
}

func parseEasyHTTP(r *http.Request) ([]byte, error) {
    return httputil.DumpRequest(r, true)
}

func parseEasyHTTPWithHost(r *http.Request, host string) ([]byte, error) {
    r.Host = host
    return httputil.DumpRequest(r, true)
}

func copyEasyHTTP(dst http.ResponseWriter, src *http.Response) (error) {
    body, err := ioutil.ReadAll(src.Body)
    if err != nil {
        return err
    }

    dstHeaders := dst.Header()
    srcHeaders := src.Header

    for header, v := range srcHeaders {
        for _, headerValue := range v {
            dstHeaders.Add(header, headerValue)
        }
    }

    dst.WriteHeader(src.StatusCode)
    dst.Write(body)
    return nil
}