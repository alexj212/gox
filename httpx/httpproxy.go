package httpx

import (
    "bytes"
    "io/ioutil"
    "net"
    "net/http"
    "net/http/httputil"
    "strconv"
    "time"
)

type ProxyOverride struct {
    Header           string
    Match            string
    Host             string
    Path             string
    ResponseModifier BytesModifier
}

type ProxyConfig struct {
    Path     string
    Host     string
    Scheme   string
    Override ProxyOverride
}

func GenerateProxy(conf ProxyConfig) http.Handler {
    proxy := &httputil.ReverseProxy{Director: func(req *http.Request) {
        originHost := conf.Host
        req.Header.Add("X-Forwarded-Host", req.Host)
        req.Header.Add("X-Origin-Host", originHost)
        req.Host = originHost
        req.URL.Host = originHost
        req.URL.Scheme = conf.Scheme // "https"

        if conf.Override.Header != "" && conf.Override.Match != "" {
            if req.Header.Get(conf.Override.Header) == conf.Override.Match {
                req.URL.Path = conf.Override.Path
            }
        }

    }, Transport: &http.Transport{
        Dial: (&net.Dialer{
            Timeout: 5 * time.Second,
        }).Dial,
    }}

    if conf.Override.ResponseModifier != nil {
        proxy.ModifyResponse = BodyRewriter(conf.Override.ResponseModifier)
    }

    return proxy
}

type BytesModifier func([]byte) []byte

func BodyRewriter(bytesModifier BytesModifier) (handler func(resp *http.Response) (err error)) {
    return func(resp *http.Response) (err error) {

        b, err := ioutil.ReadAll(resp.Body) //Read html
        if err != nil {
            return err
        }
        err = resp.Body.Close()
        if err != nil {
            return err
        }

        if bytesModifier != nil {
            b = bytesModifier(b)
        }

        body := ioutil.NopCloser(bytes.NewReader(b))
        resp.Body = body
        resp.ContentLength = int64(len(b))
        resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
        return nil
    }
}