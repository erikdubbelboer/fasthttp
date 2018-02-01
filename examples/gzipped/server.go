package main

import (
	"log"

	"github.com/erikdubbelboer/fasthttp"
)

func main() {
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
