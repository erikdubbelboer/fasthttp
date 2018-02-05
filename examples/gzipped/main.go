package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/erikdubbelboer/fasthttp"
)

func server() {
	server := fasthttp.Server{
		Name:              "Fasthttp server",
		Handler:           handler,
		ReduceMemoryUsage: true,
	}
	log.Fatal(server.ListenAndServe(":1313"))
}

func handler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html")
	ctx.Response.Header.Add("Content-Encoding", "gzip")
	if ctx.Request.Header.HasAcceptEncoding("gzip") {
		log.Println("Sending gzipped content")
		ctx.Write(
			fasthttp.AppendGzipBytes(
				nil, []byte(`<html><head><title>Compressed</title></head><body>Hello</body></html>`),
			),
		)
	} else {
		log.Println("Sending plain content")
		ctx.Write(
			[]byte(`<html><head><title>Not compressed</title></head><body>Hello</body></html>`),
		)
	}
}

func main() {
	go server()
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()

	req.Header.Add("Accept-Encoding", "gzip")
	req.SetRequestURI("http://localhost:1313")

	err := fasthttp.Do(req, res)
	if err != nil {
		panic(err)
	}
	body := res.Body()
	if b := res.Header.Peek("Content-Encoding"); len(b) > 0 {
		if bytes.Index(b, []byte("gzip")) >= 0 {
			body, err = res.BodyGunzip()
			if err != nil {
				panic(err)
			}
		}
	}
	fmt.Printf("%s\n", body)
}
