// Example static file server.
//
// Serves static files from the given directory.
// Exports various stats at /stats .
package main

import (
	"expvar"
	"flag"
	"log"
	"encoding/json"
	"github.com/erikdubbelboer/fasthttp"
	"github.com/erikdubbelboer/fasthttp/expvarhandler"
	"fmt"
	"bufio"
	"time"
	"net/http"
)

var (
	addr               = flag.String("addr", "localhost:8080", "TCP address to listen to")
	addrTLS            = flag.String("addrTLS", "", "TCP address to listen to TLS (aka SSL or HTTPS) requests. Leave empty for disabling TLS")
	byteRange          = flag.Bool("byteRange", false, "Enables byte range requests if set to true")
	certFile           = flag.String("certFile", "./ssl-cert-snakeoil.pem", "Path to TLS certificate file")
	compress           = flag.Bool("compress", false, "Enables transparent response compression if set to true")
	dir                = flag.String("dir", "/usr/share/nginx/html", "Directory to serve static files from")
	generateIndexPages = flag.Bool("generateIndexPages", true, "Whether to generate directory index pages")
	keyFile            = flag.String("keyFile", "./ssl-cert-snakeoil.key", "Path to TLS key file")
	vhost              = flag.Bool("vhost", false, "Enables virtual hosting by prepending the requested path with the requested hostname")
)

func hellloHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "<h1>Hello, world!</h1>")

	ctx.SetContentType("text/html; charset=utf8")

	// Set arbitrary headers
	ctx.Response.Header.Set("X-My-Header", "my-header-value")

	// Set cookies
	var c fasthttp.Cookie
	c.SetKey("cookie-name")
	c.SetValue("cookie-value")
	ctx.Response.Header.SetCookie(&c)
}

func jsonHandler(ctx *fasthttp.RequestCtx) {
	type Vertex struct {
		X, Y int
		Rock string
		Z    int
	}
	v1 := Vertex{1, 2, "Thay Cuong is happy", 10}
	b, err := json.Marshal(&v1)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx.Response.Header.Add("Content-Encoding", "gzip")
	ctx.SetContentType("application/json; charset=utf8")
	/* This is ok
	ctx.Write(
		fasthttp.AppendGzipBytes(
			nil, b),
	)*/
	fasthttp.WriteGzip(ctx, b)
}


//Demo streaming
func streamHandler(ctx *fasthttp.RequestCtx) {
	type Geolocation struct {
		Altitude  float64
		Latitude  float64
		Longitude float64
	}

	locations := []Geolocation{
		{-97, 37.819929, -122.478255},
		{1899, 39.096849, -120.032351},
		{2619, 37.865101, -119.538329},
		{42, 33.812092, -117.918974},
		{15, 37.77493, -122.419416},
		{2613, 67.865101, -119.538329},
		{44, 53.812092, -117.918974},
		{25, 57.77493, -122.419416},
	}

	ctx.SetContentType("application/json; charset=utf8")

	closeStream := make(chan int)
	streamReader := fasthttp.NewStreamReader(func(w *bufio.Writer) {
		for _, l := range locations {
			if err := json.NewEncoder(w).Encode(l); err != nil {
				continue
			}
			w.Flush()
			time.Sleep(1 * time.Second)
		}
		closeStream <- 0  //Signal to close stream reader
	})

	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBodyStream(streamReader, -1)
	go func() {
		select {
		case <-closeStream:
			log.Print("Stream Reader Close !")
			streamReader.Close()
		}
	}()

}
func main() {
	// Parse command-line flags.
	flag.Parse()

	// Setup FS handler
	fs := &fasthttp.FS{
		Root:               *dir,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: *generateIndexPages,
		Compress:           *compress,
		AcceptByteRange:    *byteRange,
	}
	if *vhost {
		fs.PathRewrite = fasthttp.NewVHostPathRewriter(0)
	}
	fsHandler := fs.NewRequestHandler()

	// Create RequestHandler serving server stats on /stats and files
	// on other requested paths.
	// /stats output may be filtered using regexps. For example:
	//
	//   * /stats?r=fs will show only stats (expvars) containing 'fs'
	//     in their names.
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/stats":
			expvarhandler.ExpvarHandler(ctx)
		case "/hello":
			hellloHandler(ctx)
		case "/json": //Demo JSON with gzip encoding
			jsonHandler(ctx)
		case "/stream": //Demo streaming
			streamHandler(ctx)
		default:
			fsHandler(ctx)
			updateFSCounters(ctx)
		}
	}

	// Start HTTP server.
	if len(*addr) > 0 {
		log.Printf("Starting HTTP server on %q", *addr)
		go func() {
			if err := fasthttp.ListenAndServe(*addr, requestHandler); err != nil {
				log.Fatalf("error in ListenAndServe: %s", err)
			}
		}()
	}

	// Start HTTPS server.
	if len(*addrTLS) > 0 {
		log.Printf("Starting HTTPS server on %q", *addrTLS)
		go func() {
			if err := fasthttp.ListenAndServeTLS(*addrTLS, *certFile, *keyFile, requestHandler); err != nil {
				log.Fatalf("error in ListenAndServeTLS: %s", err)
			}
		}()
	}

	log.Printf("Serving files from directory %q", *dir)
	log.Printf("See stats at http://%s/stats", *addr)

	// Wait forever.
	select {}
}

func updateFSCounters(ctx *fasthttp.RequestCtx) {
	// Increment the number of fsHandler calls.
	fsCalls.Add(1)

	// Update other stats counters
	resp := &ctx.Response
	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		fsOKResponses.Add(1)
		fsResponseBodyBytes.Add(int64(resp.Header.ContentLength()))
	case fasthttp.StatusNotModified:
		fsNotModifiedResponses.Add(1)
	case fasthttp.StatusNotFound:
		fsNotFoundResponses.Add(1)
	default:
		fsOtherResponses.Add(1)
	}
}

// Various counters - see https://golang.org/pkg/expvar/ for details.
var (
	// Counter for total number of fs calls
	fsCalls = expvar.NewInt("fsCalls")

	// Counters for various response status codes
	fsOKResponses          = expvar.NewInt("fsOKResponses")
	fsNotModifiedResponses = expvar.NewInt("fsNotModifiedResponses")
	fsNotFoundResponses    = expvar.NewInt("fsNotFoundResponses")
	fsOtherResponses       = expvar.NewInt("fsOtherResponses")

	// Total size in bytes for OK response bodies served.
	fsResponseBodyBytes = expvar.NewInt("fsResponseBodyBytes")
)
